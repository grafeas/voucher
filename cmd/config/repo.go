package config

import (
	"context"
	"net/url"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/Shopify/voucher/repository"
	"github.com/Shopify/voucher/repository/github"
)

func newRepositoryClient(org repository.Organization) (repositoryClient repository.Client) {
	orgURL, err := normalizeURL(org.URL)
	if nil != err {
		log.Errorf("failed to parse organization URL (%s): %s", org.URL, err)
		return nil
	}
	if strings.HasSuffix(orgURL.Hostname(), "github.com") {
		keyring, err := getRepositoryKeyRing()
		if nil != err {
			log.Errorf("failed to get github token: %s", err)
			return nil
		}
		for domain, token := range keyring {
			repoURL, err := normalizeURL(domain)
			if isOrgURL(orgURL, repoURL) && err == nil {
				repositoryClient, err = github.NewClient(context.Background(), &token)
				if nil != err {
					log.Errorf("failed to connect to github: %s", err)
					return nil
				}
				return repositoryClient
			}
		}
	}
	return nil
}

func normalizeURL(urlStr string) (*url.URL, error) {
	if !strings.HasPrefix(urlStr, "http://") && !strings.HasPrefix(urlStr, "https://") {
		return url.Parse("https://" + urlStr)
	}
	return url.Parse(urlStr)
}

func isOrgURL(orgURL, repoURL *url.URL) bool {
	if orgURL.Hostname() != repoURL.Hostname() {
		return false
	}
	return strings.HasPrefix(repoURL.EscapedPath(), orgURL.EscapedPath())
}
