package repository

// Organization contains repository information pertaining to an organization
type Organization struct {
	ID   string
	Name string
	URL  string
}

// CreateNewOrganization returns a new Organization object
func CreateNewOrganization(id string, name string, url string) Organization {
	return Organization{
		ID:   id,
		Name: name,
		URL:  url,
	}
}

// CommitInfo contains information pertaining to the validity of a commit
type CommitInfo struct {
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

// CreateNewCheck returns a new Check object
func CreateNewCheck(appName string, appURL string, status string, conclusion string) Check {
	return Check{
		App: App{
			Name: appName,
			URL:  appURL,
		},
		Status:     status,
		Conclusion: conclusion,
	}
}

// CreateNewCommit returns a new Commit object
func CreateNewCommit(commitURL string) Commit {
	return Commit{
		URL: commitURL,
	}
}

// CreateNewPullRequest returns a new PullRequest object
func CreateNewPullRequest(baseBranchName string, headBranchName string, isMerged bool) PullRequest {
	return PullRequest{
		BaseBranchName: baseBranchName,
		HeadBranchName: headBranchName,
		IsMerged:       isMerged,
	}
}

// CreateNewCommitInfo returns a new CommitInfo object
func CreateNewCommitInfo(commitURL string, checks []Check, commitStatus string, isSigned bool, associatedPullRequests []PullRequest) CommitInfo {
	return CommitInfo{
		URL:                    commitURL,
		Checks:                 checks,
		Status:                 commitStatus,
		IsSigned:               isSigned,
		AssociatedPullRequests: associatedPullRequests,
	}
}

// DefaultBranch contains the information related to the repository's default branch
type DefaultBranch struct {
	Name    string
	Commits []Commit
}

// CreateNewDefaultBranch returns a new DefaultBranch object
func CreateNewDefaultBranch(name string, commits []Commit) DefaultBranch {
	return DefaultBranch{
		Name:    name,
		Commits: commits,
	}
}

// Commit contains information pertaining to a commit
type Commit struct {
	URL string
}

// PullRequest contains information pertaining to a pull request
type PullRequest struct {
	BaseBranchName string
	HeadBranchName string
	IsMerged       bool
}
