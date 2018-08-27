package main

import (
	"context"
	"fmt"
	"time"

	"github.com/Shopify/voucher"
	"github.com/Shopify/voucher/auth/google"
	"github.com/Shopify/voucher/client"
	"github.com/docker/distribution/reference"
)

const timeout = 120 * time.Second

func lookupAndAttest(hostname, check, image string) error {
	var ok bool
	var namedRef reference.Named

	ref, err := reference.Parse(image)
	if nil != err {
		return fmt.Errorf("parsing image reference failed: %s", err)
	}

	if namedRef, ok = ref.(reference.Named); !ok {
		return fmt.Errorf("couldn't get named version of reference: %s", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	voucherClient, err := voucher.AuthToClient(ctx, google.NewAuth(), namedRef)
	if nil != err {
		return fmt.Errorf("creating authenticated client failed: %s", err)
	}

	canonicalRef, err := getCanonicalReference(voucherClient, namedRef)
	if nil != err {
		return fmt.Errorf("getting image digest failed: %s", err)
	}

	fmt.Printf(" - Attesting image: %s\n", canonicalRef.String())

	voucherResp, err := client.SignImage(hostname, canonicalRef, check)
	if nil != err {
		return fmt.Errorf("signing image failed: %s", err)
	}

	fmt.Println(formatResponse(&voucherResp))

	return nil
}
