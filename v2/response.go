package voucher

import "github.com/docker/distribution/reference"

// Response describes the response from a Check call.
type Response struct {
	Image   string        `json:"image"`
	Success bool          `json:"success"`
	Results []CheckResult `json:"results"`
}

// NewResponse creates a new Response for the passed ImageData,
// with the passed results.
func NewResponse(reference reference.Reference, results []CheckResult) (checkResponse Response) {
	checkResponse.Image = reference.String()
	checkResponse.Results = results
	checkResponse.Success = true

	for _, check := range checkResponse.Results {
		if !check.Success {
			checkResponse.Success = false
			break
		}
	}

	return checkResponse
}
