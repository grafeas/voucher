package objects

import (
	"github.com/antihax/optional"
)

// ListOpts based on
// ListNotesOpts https://github.com/grafeas/client-go/blob/39fa98b49d38de3942716c0f58f3505012415470/0.1.0/api_grafeas_v1_beta1.go#L1051
// ListNoteOccurrencesOpts https://github.com/grafeas/client-go/blob/39fa98b49d38de3942716c0f58f3505012415470/0.1.0/api_grafeas_v1_beta1.go#L943
type ListOpts struct {
	Filter    optional.String //not implemented for grafeas os
	PageSize  optional.Int32
	PageToken optional.String
}

// ListNotesResponse based on
// https://github.com/grafeas/client-go/blob/master/0.1.0/model_v1beta1_list_notes_response.go
type ListNotesResponse struct {
	Notes         []Note `json:"notes,omitempty"`
	NextPageToken string `json:"nextPageToken,omitempty"`
}

// ListOccurrencesResponse based on
// https://github.com/grafeas/client-go/blob/master/0.1.0/model_v1beta1_list_note_occurrences_response.go
type ListOccurrencesResponse struct {
	Occurrences   []Occurrence `json:"occurrences,omitempty"`
	NextPageToken string       `json:"nextPageToken,omitempty"`
}
