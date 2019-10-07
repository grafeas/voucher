package github

import (
	"errors"
	"regexp"
	"net/url"
	"path"
	"github.com/Shopify/voucher"
)

const (
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
func GetCommitURL(b *voucher.BuildDetail) (string, error) {
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

func ParseGithubURL(gitURL string) (*RepositoryMetadata, error) {
	re, err := regexp.Compile(GithubRegex)
	if err != nil {
		return nil, err
	}
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
	return repo, err
}
