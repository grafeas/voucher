package repository

import (
	"net/url"
	"path"
	"regexp"
)

const (
	Protocol      = "((?P<protocol>https?|git)(?:://|@))?"
	VCSName       = "(?P<vcs>[^/:]+)"
	OrgName       = "(?P<org>[^/.]+)"
	RepoName      = "(?P<repo>[^/.]+)"
	RepoExtension = "(?:.git)?"

	VCSRegex  = "^" + Protocol + VCSName + "/?$"
	OrgRegex  = "^" + Protocol + VCSName + "[/:]" + OrgName + "/?$"
	RepoRegex = "^" + Protocol + VCSName + "[/:]" + OrgName + "/" + RepoName + RepoExtension + "$"
)

// Organization contains repository information pertaining to an organization
type Organization struct {
	Alias string
	VCS   string
	Name  string
}

func NewOrganization(alias string, url string) *Organization {
	for _, regex := range []string{
		RepoRegex,
		OrgRegex,
		VCSRegex,
	} {
		if matched, err := regexp.MatchString(regex, url); nil == err && matched {
			match := getMatchGroups(regex, url)

			if "" == alias {
				alias = match["org"]
			}

			return &Organization{
				Alias: alias,
				VCS:   match["vcs"],
				Name:  match["org"],
			}
		}
	}
	return nil
}

const (
	CommitStatusError    = "ERROR"
	CommitStatusExpected = "EXPECTED"
	CommitStatusFAilure  = "FAILURE"
	CommitStatusPending  = "PENDING"
	CommitStatusSuccess  = "SUCCESS"
)

// Commit contains information pertaining to the validity of a commit
type Commit struct {
	URL                    string
	Checks                 []Check
	Status                 string
	IsSigned               bool
	AssociatedPullRequests []PullRequest
}

// Check is a collection of the check runs created by a single CI/CD App for a specific commit
type Check struct {
	Status     string
	Conclusion string
}

// NewCheck returns a new Check object
func NewCheck(status string, conclusion string) Check {
	return Check{
		Status:     status,
		Conclusion: conclusion,
	}
}

// App contains the relevant information associated with a CI/CD app
type App struct {
	Name string
	URL  string
}

// NewCommit returns a new Commit object
func NewCommit(commitURL string, checks []Check, commitStatus string, isSigned bool, associatedPullRequests []PullRequest) Commit {
	return Commit{
		URL:                    commitURL,
		Checks:                 checks,
		Status:                 commitStatus,
		IsSigned:               isSigned,
		AssociatedPullRequests: associatedPullRequests,
	}
}

// Branch contains the information related to the repository's default branch
type Branch struct {
	Name       string
	CommitRefs []CommitRef
}

// NewBranch returns a new Branch object
func NewBranch(name string, commits []CommitRef) Branch {
	return Branch{
		Name:       name,
		CommitRefs: commits,
	}
}

// CommitRef contains a URL referencing a commit
type CommitRef struct {
	URL string
}

// NewCommitRef returns a new CommitRef object
func NewCommitRef(commitURL string) CommitRef {
	return CommitRef{
		URL: commitURL,
	}
}

// PullRequest contains information pertaining to a pull request
type PullRequest struct {
	BaseBranchName       string
	HeadBranchName       string
	IsMerged             bool
	MergeCommit          CommitRef
	HasRequiredApprovals bool
}

// NewPullRequest returns a new PullRequest object
func NewPullRequest(baseBranchName string, headBranchName string, isMerged bool, mergeCommit CommitRef, hasRequiredApprovals bool) PullRequest {
	return PullRequest{
		BaseBranchName:       baseBranchName,
		HeadBranchName:       headBranchName,
		IsMerged:             isMerged,
		MergeCommit:          mergeCommit,
		HasRequiredApprovals: hasRequiredApprovals,
	}
}

// Metadata describes the top level metadata information about a repo
// that one can get from the gitUrl
type Metadata struct {
	VCS          string `json:"vcs"`
	Organization string `json:"organization"`
	Name         string `json:"name"`
}

func NewRepositoryMetadata(url string) *Metadata {
	for _, regex := range []string{
		RepoRegex,
		OrgRegex,
		VCSRegex,
	} {
		if matched, err := regexp.MatchString(regex, url); nil == err && matched {
			match := getMatchGroups(regex, url)
			return &Metadata{
				VCS:          match["vcs"],
				Organization: match["org"],
				Name:         match["repo"],
			}
		}
	}
	return nil
}

func (metadata *Metadata) String() string {
	scheme, _ := url.Parse("https://")
	scheme.Path = path.Join(
		scheme.Path,
		metadata.VCS,
		metadata.Organization,
		metadata.Name,
	)
	return scheme.String()
}

func getMatchGroups(regEx, url string) (paramsMap map[string]string) {
	var compRegEx = regexp.MustCompile(regEx)
	match := compRegEx.FindStringSubmatch(url)

	paramsMap = make(map[string]string)
	for i, name := range compRegEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return
}
