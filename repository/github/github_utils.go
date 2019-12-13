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

// RepositoryMetadata describes the top level metadata information about a repo
// that one can get from the gitUrl
type RepositoryMetadata struct {
	Vcs          string `json:"vcs"`
	Organization string `json:"organization"`
	Name         string `json:"name"`
}

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
// RepositryMetadata with the information contained in that URL.
func ParseGithubURL(gitURL string) (*RepositoryMetadata, error) {
	re := regexp.MustCompile(GithubRegex)
	repoDetail := re.FindStringSubmatch(gitURL)

	// ensure regex matches
	if len(repoDetail) != 4 {
		return nil, errors.New("unable to parse github url")
	}
	repo := &RepositoryMetadata{
		Vcs:          repoDetail[1],
		Organization: repoDetail[2],
		Name:         repoDetail[3],
	}
	return repo, nil
}
