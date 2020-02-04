package repository

// Organization contains repository information pertaining to an organization
type Organization struct {
	Name string
	URL  string
}

// NewOrganization returns a new Organization object
func NewOrganization(name string, url string) Organization {
	return Organization{
		Name: name,
		URL:  url,
	}
}

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
	App        App
	Status     string
	Conclusion string
}

// App contains the relevant information associated with a CI/CD app
type App struct {
	Name string
	URL  string
}

// NewCheck returns a new Check object
func NewCheck(appName string, appURL string, status string, conclusion string) Check {
	return Check{
		App: App{
			Name: appName,
			URL:  appURL,
		},
		Status:     status,
		Conclusion: conclusion,
	}
}

// NewCommitRef returns a new CommitRef object
func NewCommitRef(commitURL string) CommitRef {
	return CommitRef{
		URL: commitURL,
	}
}

// NewPullRequest returns a new PullRequest object
func NewPullRequest(baseBranchName string, headBranchName string, isMerged bool) PullRequest {
	return PullRequest{
		BaseBranchName: baseBranchName,
		HeadBranchName: headBranchName,
		IsMerged:       isMerged,
	}
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

// PullRequest contains information pertaining to a pull request
type PullRequest struct {
	BaseBranchName string
	HeadBranchName string
	IsMerged       bool
}

// Metadata describes the top level metadata information about a repo
// that one can get from the gitUrl
type Metadata struct {
	Vcs          string `json:"vcs"`
	Organization string `json:"organization"`
	Name         string `json:"name"`
}

// NewOrganization returns a new Organization object
func NewRepositoryMetadata(vcs, org, name string) Metadata {
	return Metadata{
		Vcs:          vcs,
		Organization: org,
		Name:         name,
	}
}
