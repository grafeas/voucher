package github

import (
	"net/url"
	"path"

	"github.com/grafeas/voucher/repository"
)

// GetCommitURL generates a commit url from the build metadata
func GetCommitURL(b *repository.BuildDetail) (string, error) {
	repoMeta := repository.NewRepositoryMetadata(b.RepositoryURL)
	scheme, err := url.Parse(repoMeta.String())
	if err != nil {
		return "", err
	}

	scheme.Path = path.Join(
		scheme.Path,
		"commit",
		b.Commit,
	)
	return scheme.String(), nil
}
