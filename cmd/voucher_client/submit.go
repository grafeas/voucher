package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/docker/distribution/reference"

	"github.com/grafeas/voucher"
)

var errImageCheckFailed = errors.New("image failed to pass required check(s)")

// check checks the passed image to the voucher server.
func check(ctx context.Context, client voucher.Interface, check string, canonicalRef reference.Canonical) error {
	fmt.Printf("Submitting image to Voucher: %s\n", canonicalRef.String())

	voucherResp, err := client.Check(ctx, check, canonicalRef)
	if nil != err {
		return fmt.Errorf("signing image failed: %s", err)
	}

	fmt.Println(formatResponse(&voucherResp))

	if !voucherResp.Success {
		return errImageCheckFailed
	}

	return nil
}

// LookupAndCheck looks up the passed image, and checks it with the Voucher
// server.
func LookupAndCheck(args []string) {
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

	err = check(ctx, client, getCheck(), canonicalRef)
	if nil != err {
		errorf("checking image with voucher failed: %s", err)
		os.Exit(1)
	}
}
