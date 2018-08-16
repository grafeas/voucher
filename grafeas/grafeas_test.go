package grafeas

import (
	"testing"

	"github.com/Shopify/voucher"
)

const testHostname = "gcr.io/alpine/alpine@sha256:297524b7375fbf09b3784f0bbd9cb2505700dd05e03ce5f5e6d262bf2f5ac51c"

const testResourceAddress = "resourceUrl=\"https://" + testHostname + "\""

func TestGrafeasHelperFunctions(t *testing.T) {
	imageData, err := voucher.NewImageData(testHostname)
	if nil != err {
		t.Fatal(err)
	}

	if testResourceAddress != resourceURL(imageData) {
		t.Errorf("Expected:\t%s\nGot:\t%s\n", testResourceAddress, resourceURL(imageData))
	}

}
