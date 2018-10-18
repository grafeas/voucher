package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Shopify/voucher"
	"github.com/Shopify/voucher/auth/google"
	"github.com/Shopify/voucher/client"
	"github.com/docker/distribution/reference"
)

// lookupCanonical looks up the canonical version of the passed image path.
func lookupCanonical(ctx context.Context, image string) (reference.Canonical, error) {
	var ok bool
	var namedRef reference.Named

	ref, err := reference.Parse(image)
	if nil != err {
		return nil, fmt.Errorf("parsing image reference failed: %s", err)
	}

	if namedRef, ok = ref.(reference.Named); !ok {
		return nil, fmt.Errorf("couldn't get named version of reference: %s", err)
	}

	voucherClient, err := voucher.AuthToClient(ctx, google.NewAuth(), namedRef)
	if nil != err {
		return nil, fmt.Errorf("creating authenticated client failed: %s", err)
	}

	canonicalRef, err := getCanonicalReference(voucherClient, namedRef)
	if nil != err {
		err = fmt.Errorf("getting image digest failed: %s", err)
	}

	return canonicalRef, err
}

// submit submits the passed image to the voucher server.
func submit(client *client.VoucherClient, check string, canonicalRef reference.Canonical) error {
	fmt.Printf("Submitting image to Voucher: %s\n", canonicalRef.String())

	voucherResp, err := client.Check(check, canonicalRef)
	if nil != err {
		return fmt.Errorf("signing image failed: %s", err)
	}

	fmt.Println(formatResponse(&voucherResp))

	return nil
}

// LookupAndSubmit looks up the passed image, and submits it with the Voucher server.
func LookupAndSubmit(args []string) {
	var err error

	client, err := getVoucherClient()
	if nil != err {
		errorf("creating client failed: %s", err)
		os.Exit(1)
	}

	ctx, cancel := newContext()
	defer cancel()

	canonicalRef, err := lookupCanonical(ctx, args[0])
	if nil != err {
		errorf("getting canonical reference failed: %s", err)
		os.Exit(1)
	}

	err = submit(client, getCheck(), canonicalRef)
	if nil != err {
		errorf("submitting image to voucher failed: %s", err)
		os.Exit(1)
	}
}
