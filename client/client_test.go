package client

import (
	"strings"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	_, err := NewClient("", 50*time.Second)
	if errNoHost != err {
		t.Fatal("should have been a no-host error, is actually: ", err)
	}

	client, err := NewClient("localhost", 50*time.Second)
	if nil != err {
		t.Fatal("failed to create client: ", err)
	}

	if nil == client.Hostname {
		t.Fatal("client hostname URL is nil")
	}

	if "https://localhost" != client.Hostname.String() {
		t.Errorf(client.Hostname.String())
	}

	_, err = NewClient(":localhost", 50*time.Second)
	if !strings.HasPrefix(err.Error(), "could not parse voucher hostname") {
		t.Fatal("failed to create client: ", err)
	}

}

func TestVoucherURL(t *testing.T) {
	client, err := NewClient("localhost", 50*time.Second)
	if nil != err {
		t.Fatal("failed to create client: ", err)
	}

	allTestURL := toVoucherURL(client.Hostname, "all")
	if "https://localhost/all" != allTestURL {
		t.Errorf("url is incorrect, should be \"%s\" instead of \"%s\"", "https://localhost/all", allTestURL)
	}

	allEmptyURL := toVoucherURL(nil, "all")
	if "/all" != allEmptyURL {
		t.Errorf("url is incorrect, should be \"%s\" instead of \"%s\"", "/all", allEmptyURL)
	}
}
