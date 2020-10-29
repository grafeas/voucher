package github

import (
	"fmt"
	"net/url"

	"github.com/grafeas/voucher/repository"
	"github.com/shurcooL/githubv4"
)

// createNewGitHubV4URI creates a new URI object used in GitHub's GraphQL queries
func createNewGitHubV4URI(uri string) (*githubv4.URI, error) {
	parsedURI, err := url.Parse(uri)
	if err != nil {
		return nil, fmt.Errorf("URI %s could not be parsed. Error: %s", uri, err)
	}

	return &githubv4.URI{
		URL: parsedURI,
	}, nil
}

// typeMismatchError represents a type mismatch between objects
type typeMismatchError struct {
	expectedType string
	actualType   string
}

func (t *typeMismatchError) Error() string {
	return fmt.Sprintf("type mismatch found. Expected: %s, Actual: %s", t.expectedType, t.actualType)
}

// newTypeMismatchError creates a new TypeMismatchError
func newTypeMismatchError(expected string, actual interface{}) error {
	actualType := fmt.Sprintf("%T", actual)
	return &typeMismatchError{
		expectedType: expected,
		actualType:   actualType,
	}
}

func (s checkConclusionState) isValidCheckConclusionState() bool {
	const (
		actionRequiredString = "ACTION_REQUIRED"
		cancelledString      = "CANCELLED"
		failureString        = "FAILURE"
		neutralString        = "NEUTRAL"
		successString        = "SUCCESS"
		timedOutString       = "TIMED_OUT"
	)
	switch s {
	case "": // checkConclusionState can be inconclusive/empty
		return true
	case actionRequiredString, cancelledString, failureString, neutralString, successString, timedOutString:
		return true
	default:
		return false
	}
}

func (s checkStatusState) isValidCheckStatusState() bool {
	const (
		completedString  = "COMPLETED"
		inProgressString = "IN_PROGRESS"
		queuedString     = "QUEUED"
		requestedString  = "REQUESTED"
	)
	switch s {
	case "": // checkStatusState can be inconclusive/empty
		return true
	case completedString, inProgressString, queuedString, requestedString:
		return true
	default:
		return false
	}
}

func (s statusState) isValidStatusState() bool {
	switch s {
	case "": // statusState can be inconclusive/empty
		return true
	case repository.CommitStatusError, repository.CommitStatusExpected, repository.CommitStatusFAilure,
		repository.CommitStatusPending, repository.CommitStatusSuccess:
		return true
	default:
		return false
	}
}

const (
	pullRequestReviewApproved         = "APPROVED"
	pullRequestReviewChangesRequested = "CHANGES_REQUESTED"
	pullRequestReviewCommented        = "COMMENTED"
	pullRequestReviewDismissed        = "DISMISSED"
	pullRequestReviewPending          = "PENDING"
)

func (s pullRequestReviewState) isValidPullRequestReviewState() bool {
	switch s {
	case "": // pullRequestReviewState can be inconclusive/empty
		return true
	case pullRequestReviewApproved, pullRequestReviewChangesRequested, pullRequestReviewCommented, pullRequestReviewDismissed, pullRequestReviewPending:
		return true
	default:
		return false
	}
}
