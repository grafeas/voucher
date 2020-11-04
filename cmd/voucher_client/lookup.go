package main

import (
	"context"
	"fmt"

	"github.com/docker/distribution/reference"

	"github.com/grafeas/voucher"
	"github.com/grafeas/voucher/auth/google"
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
