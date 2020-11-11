package main

import (
	"context"
	"fmt"
	"os"

	"github.com/docker/distribution/reference"

	"github.com/grafeas/voucher"
)

// verifyImage submits the passed image to the voucher server for verification.
func verifyImage(ctx context.Context, client voucher.Interface, check string, canonicalRef reference.Canonical) error {
	fmt.Printf("Verifying image with Voucher: %s\n", canonicalRef.String())

	voucherResp, err := client.Verify(ctx, check, canonicalRef)
	if nil != err {
		return fmt.Errorf("verifying image failed: %s", err)
	}

	fmt.Println(formatResponse(&voucherResp))

	if !voucherResp.Success {
		return errImageCheckFailed
	}

	return nil
}

// LookupAndVerify looks up the passed image, and submits it with the Voucher server.
func LookupAndVerify(args []string) {
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

	err = verifyImage(ctx, client, getCheck(), canonicalRef)
	if nil != err {
		errorf("verifying image with voucher failed: %s", err)
		os.Exit(1)
	}
}
