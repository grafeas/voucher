package github

import (
	"errors"
	"net/url"
	"path"
	"regexp"

	"github.com/Shopify/voucher/repository"
)

const (
	// GithubRegex matches typical Github URLs.
	GithubRegex = "^(?:https?|git)(?:://|@)(?P<vcs>[^/:]+)[/:](?P<org>[^/.]+)/(?P<name>[^/.]+)(?:.git)?"
)

// GetCommitURL generates a commit url from the build metadata
func GetCommitURL(b *repository.BuildDetail) (string, error) {
	repoMeta, err := ParseGithubURL(b.RepositoryURL)
	if err != nil {
		return "", errors.New("Error parsing github commit url")
	}
	scheme, err := url.Parse("https://")
	if err != nil {
		return "", err
	}

	scheme.Path = path.Join(
		scheme.Path,
		repoMeta.Vcs,
		repoMeta.Organization,
		repoMeta.Name,
		"commit",
		b.Commit,
	)
	return scheme.String(), nil
}

// GetRepositoryURL generates a repository url from the build metadata
func GetRepositoryURL(b *repository.BuildDetail) (string, error) {
	repoMeta, err := ParseGithubURL(b.RepositoryURL)
	if err != nil {
		return "", errors.New("Error parsing github commit url")
	}
	scheme, err := url.Parse("https://")
	if err != nil {
		return "", err
	}

	scheme.Path = path.Join(
		scheme.Path,
		repoMeta.Vcs,
		repoMeta.Organization,
		repoMeta.Name,
	)
	return scheme.String(), nil
}

// ParseGithubURL parses the passed string as a Github URL and returns a
// RepositoryMetadata with the information contained in that URL.
func ParseGithubURL(gitURL string) (*repository.Metadata, error) {
	re := regexp.MustCompile(GithubRegex)
	repoDetail := re.FindStringSubmatch(gitURL)

	// ensure regex matches
	if len(repoDetail) != 4 {
		return nil, errors.New("unable to parse github url")
	}
	repo := repository.NewRepositoryMetadata(repoDetail[1], repoDetail[2], repoDetail[3])
	return &repo, nil
}
