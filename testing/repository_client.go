package vtesting

import (
	"context"
	"regexp"

	"github.com/Shopify/voucher/repository"
)

const (
	repoRegex   = "^(?:https?|git)(?:://|@)(?P<vcs>[^/:]+)[/:](?P<org>[^/.]+)/(?P<name>[^/.]+)(?:.git)?"
	commitRegex = "^(?:https?|git)(?:://|@)(?P<vcs>[^/:]+)[/:](?P<org>[^/.]+)/(?P<name>[^/.]+)/(?:commit)/(?P<hash>[^/.]+)"
)

type repositoryclient struct {
	repos map[string]Repository
}

//NewClient creates a new repository client
func NewClient() (RepositoryClient, error) {
	return &repositoryclient{}, nil
}

//AddRepository adds to repository to the repository client
func (c *repositoryclient) AddRepository(org repository.Organization, name string, commits map[string]repository.Commit, branches map[string]repository.Branch) {
	if c.repos == nil {
		c.repos = make(map[string]Repository)
	}
	c.repos[org.Name] = Repository{org, name, commits, branches}
}

//GetCommitInfo creates commit info with provided
func (c *repositoryclient) GetCommit(ctx context.Context, details repository.BuildDetail) (repository.Commit, error) {
	re := regexp.MustCompile(repoRegex)
	res := re.FindStringSubmatch(details.RepositoryURL)
	orgName := res[2]
	repo := c.repos[orgName]
	commit := repo.Commits[details.Commit]
	return commit, nil
}

//GetOrganization creates organization
func (c *repositoryclient) GetOrganization(ctx context.Context, details repository.BuildDetail) (repository.Organization, error) {
	var org repository.Organization
	re := regexp.MustCompile(repoRegex)
	orgName := re.FindStringSubmatch(details.RepositoryURL)[2]
	if repo, exists := c.repos[orgName]; exists {
		org = repo.Org
	}
	return org, nil
}

func (c *repositoryclient) GetDefaultBranch(ctx context.Context, details repository.BuildDetail) (repository.Branch, error) {
	var branch repository.Branch
	re := regexp.MustCompile(repoRegex)
	reCommit := regexp.MustCompile(commitRegex)
	orgName := re.FindStringSubmatch(details.RepositoryURL)[2]

	branches := c.repos[orgName].Branches
	for _, br := range branches {
		commits := br.CommitRefs
		for _, c := range commits {
			hash := reCommit.FindStringSubmatch(c.URL)[4]
			if hash == details.Commit {
				branch = br
				break
			}
		}
	}
	return branch, nil
}

func (c *repositoryclient) GetBranch(ctx context.Context, details repository.BuildDetail, name string) (repository.Branch, error) {
	var branch repository.Branch
	re := regexp.MustCompile(repoRegex)
	reCommit := regexp.MustCompile(commitRegex)
	orgName := re.FindStringSubmatch(details.RepositoryURL)[2]

	master := c.repos[orgName].Branches[name]
	commits := master.CommitRefs
	for _, c := range commits {
		hash := reCommit.FindStringSubmatch(c.URL)[4]
		if hash == details.Commit {
			branch = master
			break
		}
	}
	return branch, nil
}
