package voucher

// Response describes the response from a Check call.
type Response struct {
	Image   string        `json:"image"`
	Project string        `json:"project"`
	Success bool          `json:"success"`
	Results []CheckResult `json:"results"`
}

// NewResponse creates a new Response for the passed ImageData,
// with the passed results.
func NewResponse(imageData ImageData, results []CheckResult) (checkResponse Response) {
	checkResponse.Image = imageData.String()
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
