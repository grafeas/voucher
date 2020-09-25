package rest

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/Shopify/voucher/grafeas/rest/objects"
)

//GrafeasAPIService vars
var (
	ErrTimeout = errors.New("timeout error when getting REST data")
	timeout    = time.Minute
)

//GrafeasAPIService is the interface for communicating with Grafeas
type GrafeasAPIService interface {
	CreateOccurrence(context.Context, string, objects.Occurrence) (objects.Occurrence, error)
	BatchCreateOccurrences(context.Context, string, []objects.Occurrence) ([]objects.Occurrence, error)
	BatchCreateNotes(context.Context, string, map[string]objects.Note) ([]objects.Note, error)
	ListNotes(context.Context, string, *objects.ListOpts) (objects.ListNotesResponse, error)
	ListOccurrences(context.Context, string, *objects.ListOpts) (objects.ListOccurrencesResponse, error)
	//ListNoteOccurrences(context.Context, string, string, *objects.ListOpts) (objects.ListOccurrencesResponse, error)
}

//GrafeasAPIServiceImpl for REST calls
type GrafeasAPIServiceImpl struct {
	basePath    string
	versionPath string
	client      *http.Client
}

//NewGrafeasAPIService creates new GrafeasAPIService
func NewGrafeasAPIService(basePath, versionPath string) GrafeasAPIServiceImpl {
	return GrafeasAPIServiceImpl{
		basePath:    basePath,
		versionPath: versionPath,
		client: &http.Client{
			Timeout: timeout,
		},
	}
}

//CreateOccurrence https://github.com/grafeas/client-go/blob/39fa98b49d38de3942716c0f58f3505012415470/0.1.0/api_grafeas_v1_beta1.go#L310
func (g GrafeasAPIServiceImpl) CreateOccurrence(ctx context.Context, parent string, occurrence objects.Occurrence) (objects.Occurrence, error) {
	urlPath, err := g.buildURL(parent, "/occurrences", nil)
	if err != nil {
		return objects.Occurrence{}, err
	}
	res, err := json.Marshal(occurrence)
	if err != nil {
		return objects.Occurrence{}, err
	}
	resp, err := g.httpCall(urlPath, res, http.MethodPost)
	if err != nil {
		return objects.Occurrence{}, err
	}
	occ := objects.Occurrence{}
	err = json.Unmarshal(resp, &occ)
	if err != nil {
		return objects.Occurrence{}, err
	}
	return occ, nil
}

//BatchCreateOccurrences https://github.com/grafeas/client-go/blob/39fa98b49d38de3942716c0f58f3505012415470/0.1.0/api_grafeas_v1_beta1.go#L128
func (g GrafeasAPIServiceImpl) BatchCreateOccurrences(ctx context.Context, parent string, occs []objects.Occurrence) ([]objects.Occurrence, error) {
	occsReq := objects.BatchCreateOccurrencesRequest{
		Parent:      parent,
		Occurrences: occs,
	}
	urlPath, err := g.buildURL(parent, "/occurrences:batchCreate", nil)
	if err != nil {
		return []objects.Occurrence{}, err
	}
	res, err := json.Marshal(occsReq)
	if err != nil {
		return []objects.Occurrence{}, err
	}
	resp, err := g.httpCall(urlPath, res, http.MethodPost)
	if err != nil {
		return []objects.Occurrence{}, err
	}
	occResp := objects.BatchCreateOccurrencesResponse{}
	err = json.Unmarshal(resp, &occResp)
	if err != nil {
		return []objects.Occurrence{}, err
	}
	return occResp.Occurrences, nil
}

//BatchCreateNotes https://github.com/grafeas/client-go/blob/39fa98b49d38de3942716c0f58f3505012415470/0.1.0/api_grafeas_v1_beta1.go#L37
func (g GrafeasAPIServiceImpl) BatchCreateNotes(ctx context.Context, parent string, notes map[string]objects.Note) ([]objects.Note, error) {
	notesReq := objects.BatchCreateNotesRequest{
		Parent: parent,
		Notes:  notes,
	}
	urlPath, err := g.buildURL(parent, "/notes:batchCreate", nil)
	if err != nil {
		return []objects.Note{}, err
	}
	res, err := json.Marshal(notesReq)
	if err != nil {
		return []objects.Note{}, err
	}
	resp, err := g.httpCall(urlPath, res, http.MethodPost)
	if err != nil {
		return []objects.Note{}, err
	}
	notesResp := objects.BatchCreateNotesResponse{}
	err = json.Unmarshal(resp, &notesResp)
	if err != nil {
		return []objects.Note{}, err
	}
	return notesResp.Notes, nil
}

//ListNotes https://github.com/grafeas/client-go/blob/39fa98b49d38de3942716c0f58f3505012415470/0.1.0/api_grafeas_v1_beta1.go#L1057
func (g GrafeasAPIServiceImpl) ListNotes(ctx context.Context, parent string, optsNotes *objects.ListOpts) (objects.ListNotesResponse, error) {
	urlPath, err := g.buildURL(parent, "/notes", optsNotes)
	if err != nil {
		return objects.ListNotesResponse{}, err
	}
	resp, err := g.httpCall(urlPath, nil, http.MethodGet)
	if err != nil {
		return objects.ListNotesResponse{}, err
	}
	notesResp := objects.ListNotesResponse{}
	err = json.Unmarshal(resp, &notesResp)
	if err != nil {
		return objects.ListNotesResponse{}, err
	}
	return notesResp, nil
}

//ListOccurrences https://github.com/grafeas/client-go/blob/39fa98b49d38de3942716c0f58f3505012415470/0.1.0/api_grafeas_v1_beta1.go#L1165
func (g GrafeasAPIServiceImpl) ListOccurrences(ctx context.Context, parent string, optsOccurrences *objects.ListOpts) (objects.ListOccurrencesResponse, error) {
	urlPath, err := g.buildURL(parent, "/occurrences", optsOccurrences)
	if err != nil {
		return objects.ListOccurrencesResponse{}, err
	}
	resp, err := g.httpCall(urlPath, nil, http.MethodGet)
	if err != nil {
		return objects.ListOccurrencesResponse{}, err
	}
	occResp := objects.ListOccurrencesResponse{}
	err = json.Unmarshal(resp, &occResp)
	if err != nil {
		return objects.ListOccurrencesResponse{}, err
	}
	return occResp, nil
}

//ListNoteOccurrences https://github.com/grafeas/client-go/blob/39fa98b49d38de3942716c0f58f3505012415470/0.1.0/api_grafeas_v1_beta1.go#L949
//func (g *GrafeasAPIServiceImpl) ListNoteOccurrences(ctx context.Context, parent, noteName string, optsOccurrences *objects.ListOpts) (objects.ListOccurrencesResponse, error) {
//	return g.ListOccurrences(ctx, parent+"/notes/"+noteName, optsOccurrences)
//}

func (g GrafeasAPIServiceImpl) buildURL(parent, address string, options *objects.ListOpts) (*url.URL, error) {
	path := g.basePath + g.versionPath + parent + address
	res, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	q := res.Query()
	if options != nil && options.Filter.IsSet() {
		q.Add("filter", options.Filter.Value())
	}
	if options != nil && options.PageSize.IsSet() {
		q.Add("page_size", fmt.Sprint(options.PageSize.Value()))
	}
	if options != nil && options.PageToken.IsSet() {
		q.Add("page_token", options.PageToken.Value())
	}
	res.RawQuery = q.Encode()
	return res, nil
}

func (g GrafeasAPIServiceImpl) httpCall(urlAddr *url.URL, payload []byte, method string) ([]byte, error) {
	req := http.Request{
		Method: method,
		URL:    urlAddr,
		Header: make(map[string][]string),
		Body:   ioutil.NopCloser(bytes.NewReader(payload)),
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err := g.client.Do(&req)
	if err != nil {
		urlErr := err.(*url.Error)
		if urlErr.Timeout() {
			return nil, ErrTimeout
		}
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("error getting REST data: " + strconv.Itoa(resp.StatusCode))
	}
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	return data, err
}
