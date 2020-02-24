package github

// previewSchemas contain the headers to access parts of GitHub's preview mode
var previewSchemas = []string{
	"application/vnd.github.groot-preview+json",
	"application/vnd.github.antiope-preview+json",
}

const (
	// repositoryType is one of GitHub's GraphQL schema types representing a GitHub repository
	repositoryType = "Repository"
	// organizationType is one of GitHub's GraphQL schema types representing a GitHub organization
	organizationType = "Organization"
	// commitType is one of GitHub's GraphQL schema types representing a Git commit
	commitType = "Commit"
	// pullRequestType is one of GitHub's GraphQL schema types representing a Git pull request
	pullRequestType = "PullRequest"
)

// pagination query limit
const queryPageLimit = 3

// checkConclusionState is a string that represents the state for a check suite or check run conclusion.
// checkConclusionState is a type in the GitHub v4 GraphQL Schema
type checkConclusionState string

// checkStatusState is a string that represents the state for a check suite or check run status
// checkStatusState is a type in the Github v4 GraphQL Schema
type checkStatusState string

// statusState is a string that represents the combined commit status
// statusState is a type in the Github v4 GraphQL Schema
type statusState string

// pullRequestState is a string that represents the state for a pull request review
// pullRequestReviewState is a type in the Github v4 GraphQL Schema
type pullRequestReviewState string
