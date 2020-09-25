package rest

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/Shopify/voucher/grafeas/rest/objects"
	"github.com/antihax/optional"
	"github.com/stretchr/testify/assert"
)

func TestCreateOccurrence(t *testing.T) {
	ctx := context.Background()
	data := []byte(`{"resource":{"name":"","uri":"https://gcr.io/project/image@sha256:foo","contentHash":null},"noteName":"projects/grafeastest/notes/notevuln","kind":"VULNERABILITY","remediation":"","vulnerability":{"type":"","severity":"CRITICAL","cvssScore":0,"packageIssue":[{"affectedLocation":{"cpeUri":"7","package":"a","version":{"epoch":0,"name":"v1.1.1","revision":"r","kind":"NORMAL"}},"fixedLocation":{"cpeUri":"cpe:/o:debian:debian_linux:7","package":"a","version":{"epoch":0,"name":"namestring","revision":"1","kind":"NORMAL"}},"severityName":""}],"shortDescription":"","longDescription":"","relatedUrls":[],"effectiveSeverity":"CRITICAL"}}`)
	occ := objects.Occurrence{}
	json.Unmarshal(data, &occ)
	tcs := map[string]struct {
		parent      string
		expectedOcc objects.Occurrence
		returnData  []byte
		expectedErr string
		statusCode  int
		customURL   string
	}{
		"valid input": {
			parent:      "grafeastest",
			expectedOcc: occ,
			returnData:  data,
			statusCode:  http.StatusOK,
		},
		"empty parent": {
			parent:      "",
			expectedOcc: occ,
			returnData:  data,
			statusCode:  http.StatusOK,
		},
		"json data": {
			parent:      "grafeastest",
			expectedOcc: occ,
			returnData:  data,
			statusCode:  http.StatusOK,
		},
		"invalid json": {
			parent:      "grafeastest",
			expectedOcc: occ,
			returnData:  []byte(`:{i{`),
			expectedErr: "invalid character ':' looking for beginning of value",
			statusCode:  http.StatusOK,
		},
		"err status code": {
			parent:      "grafeastest",
			expectedOcc: occ,
			statusCode:  http.StatusInternalServerError,
			expectedErr: "error getting REST data: " + strconv.Itoa(http.StatusInternalServerError),
		},
	}
	for tc, test := range tcs {
		t.Run(tc, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.statusCode)
				w.Write(test.returnData)
			}))
			defer server.Close()
			if test.customURL == "" {
				test.customURL = server.URL + "/"
			}
			client := NewGrafeasAPIService(test.customURL, "")

			resOcc, err := client.CreateOccurrence(ctx, test.parent, test.expectedOcc)

			if err == nil {
				assert.Equal(t, test.expectedOcc, resOcc)
				assert.Equal(t, nil, err)
			} else {
				assert.Equal(t, test.expectedErr, err.Error())
			}
		})
	}
}

func TestBatchCreateOccurrences(t *testing.T) {
	ctx := context.Background()
	data := []byte(`{"occurrences":[{"name":"projects/grafeasclienttest/occurrences/0841f19d-16bc-4f83-8a15-640c33fd2e4d","resource":{"name":"","uri":"https://gcr.io/project/image@sha256:foo","contentHash":null},"noteName":"projects/grafeasclienttest/notes/grafeasbuild","kind":"BUILD","remediation":"","createTime":"2020-09-21T20:13:36.120044900Z","updateTime":"2020-09-21T20:13:36.120044900Z","build":{"provenance":{"id":"id","projectId":"ShopifyId","commands":[],"builtArtifacts":[{"checksum":"sha256:71e3e78693c011e59b3fc84940f7672aeeb0a55427b6f5157bd08ab9e9ac746c","id":"some_id1","names":[]}],"createTime":"0001-01-01T00:00:00Z","startTime":"0001-01-01T00:00:00Z","endTime":"0001-01-01T00:00:00Z","creator":"Shopify","logsUri":"some_url","sourceProvenance":{"artifactStorageSourceUri":"","fileHashes":{},"context":{"git":{"url":"https://github.com/Shopify/voucher","revisionId":"q1q"},"labels":{}},"additionalContexts":[]},"triggerId":"","buildOptions":{},"builderVersion":""},"provenanceBytes":""}},{"name":"projects/grafeasclienttest/occurrences/25001d31-f407-419d-83a2-05474ce83e34","resource":{"name":"","uri":"https://gcr.io/project/image@sha256:foo","contentHash":null},"noteName":"projects/grafeasclienttest/notes/grafeasattestation","kind":"ATTESTATION","remediation":"","createTime":"2020-09-21T20:13:36.120087200Z","updateTime":"2020-09-21T20:13:36.120087200Z","attestation":{"attestation":{"pgpSignedAttestation":{"signature":"signature","contentType":"CONTENT_TYPE_UNSPECIFIED","pgpKeyId":"1234"}}}},{"name":"projects/grafeasclienttest/occurrences/332c51c9-4140-4a75-b6b2-b918aba7270d","resource":{"name":"","uri":"https://gcr.io/project/image@sha256:foo","contentHash":null},"noteName":"projects/provider_example/notes/exampleVulnerabilityNote","kind":"VULNERABILITY","remediation":"","createTime":"2020-09-22T19:13:41.283089300Z","updateTime":"2020-09-22T19:13:41.283089300Z","vulnerability":{"type":"","severity":"SEVERITY_UNSPECIFIED","cvssScore":0,"packageIssue":[{"affectedLocation":{"cpeUri":"7","package":"a","version":{"epoch":0,"name":"v1.1.1","revision":"r","kind":"NORMAL"}},"fixedLocation":{"cpeUri":"cpe:/o:debian:debian_linux:7","package":"a","version":{"epoch":0,"name":"namestring","revision":"1","kind":"NORMAL"}},"severityName":""}],"shortDescription":"","longDescription":"","relatedUrls":[],"effectiveSeverity":"SEVERITY_UNSPECIFIED"}},{"name":"projects/grafeasclienttest/occurrences/3690b2f5-cc34-4041-be95-76383cf315cb","resource":{"name":"","uri":"https://gcr.io/project/image@sha256:foo","contentHash":null},"noteName":"projects/grafeasclienttest/notes/grafeasdiscovery","kind":"VULNERABILITY","remediation":"","createTime":"2020-09-21T20:13:36.120127700Z","updateTime":"2020-09-21T20:13:36.120127700Z","vulnerability":{"type":"vulnlow","severity":"SEVERITY_UNSPECIFIED","cvssScore":0,"packageIssue":[{"affectedLocation":{"cpeUri":"uri","package":"package_test","version":{"epoch":0,"name":"v0.1.0","revision":"re","kind":"NORMAL"}},"fixedLocation":null,"severityName":""}],"shortDescription":"","longDescription":"","relatedUrls":[],"effectiveSeverity":"SEVERITY_UNSPECIFIED"}},{"name":"projects/grafeasclienttest/occurrences/7c3c4fc1-8a7e-49df-a399-8771e8778538","resource":{"name":"","uri":"https://gcr.io/project/image@sha256:foo","contentHash":null},"noteName":"projects/grafeasclienttest/notes/grafeasbuild","kind":"BUILD","remediation":"","createTime":"2020-09-21T20:13:36.120065900Z","updateTime":"2020-09-21T20:13:36.120065900Z","build":{"provenance":{"id":"provenceid","projectId":"ShopifyId2","commands":[],"builtArtifacts":[{"checksum":"sha256:71e3e78693c011e59b3fc84940f7672aeeb0a55427b6f5157bd08ab9e9ac746c","id":"some_id2","names":[]}],"createTime":"0001-01-01T00:00:00Z","startTime":"0001-01-01T00:00:00Z","endTime":"0001-01-01T00:00:00Z","creator":"Shopify","logsUri":"some_url2","sourceProvenance":{"artifactStorageSourceUri":"","fileHashes":{},"context":{"git":{"url":"https://github.com/Shopify/voucher","revisionId":"2"},"labels":{}},"additionalContexts":[]},"triggerId":"","buildOptions":{},"builderVersion":""},"provenanceBytes":""}},{"name":"projects/grafeasclienttest/occurrences/8e215cd3-2128-4005-aa54-7ee118620c9a","resource":{"name":"","uri":"https://gcr.io/project/image@sha256:foo","contentHash":null},"noteName":"projects/grafeasclienttest/notes/grafeasvulnerability","kind":"VULNERABILITY","remediation":"","createTime":"2020-09-21T20:13:36.120019600Z","updateTime":"2020-09-21T20:13:36.120019600Z","vulnerability":{"type":"vulnmedium","severity":"MINIMAL","cvssScore":0,"packageIssue":[{"affectedLocation":{"cpeUri":"uri","package":"package_test","version":{"epoch":0,"name":"v0.0.0","revision":"r","kind":"NORMAL"}},"fixedLocation":null,"severityName":""}],"shortDescription":"","longDescription":"","relatedUrls":[],"effectiveSeverity":"SEVERITY_UNSPECIFIED"}},{"name":"projects/grafeasclienttest/occurrences/9cdd69a6-ed0c-4e33-9353-ff4fd83eaf6c","resource":{"name":"","uri":"https://gcr.io/project/image@sha256:foo","contentHash":null},"noteName":"projects/grafeasclienttest/notes/grafeasdiscovery","kind":"DISCOVERY","remediation":"","createTime":"2020-09-21T20:13:36.120148600Z","updateTime":"2020-09-21T20:13:36.120148600Z","discovered":{"discovered":{"continuousAnalysis":"ACTIVE","lastAnalysisTime":null,"analysisStatus":"PENDING","analysisStatusError":null}}},{"name":"projects/grafeasclienttest/occurrences/b1c25d9f-7bc5-4ce5-b0bf-38b31bd03e6e","resource":{"name":"","uri":"https://gcr.io/project/image@sha256:foo","contentHash":null},"noteName":"projects/grafeasclienttest/notes/grafeasdiscovery","kind":"DISCOVERY","remediation":"","createTime":"2020-09-21T20:13:36.120105400Z","updateTime":"2020-09-21T20:13:36.120105400Z","discovered":{"discovered":{"continuousAnalysis":"ACTIVE","lastAnalysisTime":null,"analysisStatus":"FINISHED_SUCCESS","analysisStatusError":null}}}],"nextPageToken":""}`)
	occsR := objects.ListOccurrencesResponse{}
	json.Unmarshal(data, &occsR)

	tcs := map[string]struct {
		parent       string
		expectedOccs []objects.Occurrence
		returnData   []byte
		expectedErr  string
		statusCode   int
		customURL    string
	}{
		"valid input": {
			parent:       "grafeastest",
			expectedOccs: occsR.Occurrences,
			returnData:   data,
			statusCode:   http.StatusOK,
		},
		"empty parent": {
			parent:       "",
			expectedOccs: occsR.Occurrences,
			returnData:   data,
			statusCode:   http.StatusOK,
		},
		"json data": {
			parent:       "",
			expectedOccs: occsR.Occurrences,
			returnData:   data,
			statusCode:   http.StatusOK,
		},
		"invalid json": {
			parent:       "",
			expectedOccs: occsR.Occurrences,
			returnData:   []byte(`:{i{`),
			statusCode:   http.StatusOK,
			expectedErr:  "invalid character ':' looking for beginning of value",
		},
		"err status code": {
			parent:       "",
			expectedOccs: occsR.Occurrences,
			statusCode:   http.StatusInternalServerError,
			expectedErr:  "error getting REST data: " + strconv.Itoa(http.StatusInternalServerError),
		},
	}
	for tc, test := range tcs {
		t.Run(tc, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.statusCode)
				w.Write(test.returnData)
			}))
			defer server.Close()
			if test.customURL == "" {
				test.customURL = server.URL + "/"
			}
			client := NewGrafeasAPIService(test.customURL, "")

			resOccs, err := client.BatchCreateOccurrences(ctx, test.parent, test.expectedOccs)

			if err == nil {
				assert.Equal(t, test.expectedOccs, resOccs)
				assert.Equal(t, nil, err)
			} else {
				assert.Equal(t, test.expectedErr, err.Error())
			}
		})
	}
}

func TestBatchCreateNotes(t *testing.T) {
	ctx := context.Background()
	data := []byte(`{"notes":[{"name":"projects/grafeasclienttest/notes/grafeasattestation","shortDescription":"short","longDescription":"long","kind":"ATTESTATION","relatedUrl":[],"expirationTime":"0001-01-01T00:00:00Z","createTime":"2020-09-21T20:13:36.115326300Z","updateTime":"2020-09-21T20:13:36.115326300Z","relatedNoteNames":[],"attestationAuthority":{"hint":{"humanReadableName":"grafeasattestation"}}},{"name":"projects/grafeasclienttest/notes/grafeasbuild","shortDescription":"short","longDescription":"long","kind":"BUILD","relatedUrl":[],"expirationTime":"0001-01-01T00:00:00Z","createTime":"2020-09-21T20:13:36.115269100Z","updateTime":"2020-09-21T20:13:36.115269100Z","relatedNoteNames":[],"build":{"builderVersion":"v0.0.0","signature":null}},{"name":"projects/grafeasclienttest/notes/grafeasdiscovery","shortDescription":"short","longDescription":"long","kind":"DISCOVERY","relatedUrl":[],"expirationTime":"0001-01-01T00:00:00Z","createTime":"2020-09-21T20:13:36.115287200Z","updateTime":"2020-09-21T20:13:36.115287200Z","relatedNoteNames":[],"discovery":{"analysisKind":"VULNERABILITY"}},{"name":"projects/grafeasclienttest/notes/grafeasvulnerability","shortDescription":"short","longDescription":"long","kind":"VULNERABILITY","relatedUrl":[],"expirationTime":"0001-01-01T00:00:00Z","createTime":"2020-09-21T20:13:36.115305300Z","updateTime":"2020-09-21T20:13:36.115305300Z","relatedNoteNames":[],"vulnerability":{"cvssScore":4.3,"severity":"MINIMAL","details":[{"cpeUri":"test_url","package":"package","minAffectedVersion":null,"maxAffectedVersion":null,"severityName":"medium","description":"","fixedLocation":null,"packageType":"","isObsolete":false,"sourceUpdateTime":null}],"cvssV3":null,"windowsDetails":[],"sourceUpdateTime":null}},{"name":"projects/grafeasclienttest/notes/testNote","shortDescription":"A brief description of the note2","longDescription":"A longer description of the note","kind":"VULNERABILITY","relatedUrl":[],"expirationTime":null,"createTime":"2020-09-23T12:13:09.834847Z","updateTime":"2020-09-23T12:13:09.834847Z","relatedNoteNames":[],"vulnerability":{"cvssScore":0,"severity":"SEVERITY_UNSPECIFIED","details":[{"cpeUri":"cpe:o:debian:debian_linux:7","package":"libexempi3","minAffectedVersion":{"epoch":0,"name":"2.5.7","revision":"1","kind":"NORMAL"},"maxAffectedVersion":null,"severityName":"","description":"","fixedLocation":null,"packageType":"","isObsolete":false,"sourceUpdateTime":null}],"cvssV3":null,"windowsDetails":[],"sourceUpdateTime":null}}],"nextPageToken":""}`)
	notesR := objects.ListNotesResponse{}
	json.Unmarshal(data, &notesR)

	tcs := map[string]struct {
		parent        string
		expectedNotes []objects.Note
		returnData    []byte
		expectedErr   string
		statusCode    int
		customURL     string
	}{
		"valid input": {
			parent:        "grafeastest",
			expectedNotes: notesR.Notes,
			returnData:    data,
			statusCode:    http.StatusOK,
		},
		"empty parent": {
			parent:        "",
			expectedNotes: notesR.Notes,
			returnData:    data,
			statusCode:    http.StatusOK,
		},
		"json data": {
			parent:        "",
			expectedNotes: notesR.Notes,
			returnData:    data,
			statusCode:    http.StatusOK,
		},
		"invalid json": {
			parent:        "",
			expectedNotes: notesR.Notes,
			returnData:    []byte(`:some { string {`),
			expectedErr:   "invalid character ':' looking for beginning of value",
			statusCode:    http.StatusOK,
		},
		"err status code": {
			parent:        "",
			expectedNotes: notesR.Notes,
			statusCode:    http.StatusInternalServerError,
			expectedErr:   "error getting REST data: " + strconv.Itoa(http.StatusInternalServerError),
		},
	}
	for tc, test := range tcs {
		t.Run(tc, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.statusCode)
				w.Write(test.returnData)
			}))
			defer server.Close()
			if test.customURL == "" {
				test.customURL = server.URL + "/"
			}
			client := NewGrafeasAPIService(test.customURL, "")

			notesMap := make(map[string]objects.Note)
			for _, val := range test.expectedNotes {
				notesMap[val.Name] = val
			}

			resNotes, err := client.BatchCreateNotes(ctx, test.parent, notesMap)

			if err == nil {
				assert.Equal(t, test.expectedNotes, resNotes)
				assert.Equal(t, nil, err)
			} else {
				assert.Equal(t, test.expectedErr, err.Error())
			}
		})
	}
}

func TestListNotes(t *testing.T) {
	ctx := context.Background()
	data := []byte(`{"notes":[{"name":"projects/grafeasclienttest/notes/grafeasattestation","shortDescription":"short","longDescription":"long","kind":"ATTESTATION","relatedUrl":[],"expirationTime":"0001-01-01T00:00:00Z","createTime":"2020-09-21T20:13:36.115326300Z","updateTime":"2020-09-21T20:13:36.115326300Z","relatedNoteNames":[],"attestationAuthority":{"hint":{"humanReadableName":"grafeasattestation"}}},{"name":"projects/grafeasclienttest/notes/grafeasbuild","shortDescription":"short","longDescription":"long","kind":"BUILD","relatedUrl":[],"expirationTime":"0001-01-01T00:00:00Z","createTime":"2020-09-21T20:13:36.115269100Z","updateTime":"2020-09-21T20:13:36.115269100Z","relatedNoteNames":[],"build":{"builderVersion":"v0.0.0","signature":null}},{"name":"projects/grafeasclienttest/notes/grafeasdiscovery","shortDescription":"short","longDescription":"long","kind":"DISCOVERY","relatedUrl":[],"expirationTime":"0001-01-01T00:00:00Z","createTime":"2020-09-21T20:13:36.115287200Z","updateTime":"2020-09-21T20:13:36.115287200Z","relatedNoteNames":[],"discovery":{"analysisKind":"VULNERABILITY"}},{"name":"projects/grafeasclienttest/notes/grafeasvulnerability","shortDescription":"short","longDescription":"long","kind":"VULNERABILITY","relatedUrl":[],"expirationTime":"0001-01-01T00:00:00Z","createTime":"2020-09-21T20:13:36.115305300Z","updateTime":"2020-09-21T20:13:36.115305300Z","relatedNoteNames":[],"vulnerability":{"cvssScore":4.3,"severity":"MINIMAL","details":[{"cpeUri":"test_url","package":"package","minAffectedVersion":null,"maxAffectedVersion":null,"severityName":"medium","description":"","fixedLocation":null,"packageType":"","isObsolete":false,"sourceUpdateTime":null}],"cvssV3":null,"windowsDetails":[],"sourceUpdateTime":null}},{"name":"projects/grafeasclienttest/notes/testNote","shortDescription":"A brief description of the note2","longDescription":"A longer description of the note","kind":"VULNERABILITY","relatedUrl":[],"expirationTime":null,"createTime":"2020-09-23T12:13:09.834847Z","updateTime":"2020-09-23T12:13:09.834847Z","relatedNoteNames":[],"vulnerability":{"cvssScore":0,"severity":"SEVERITY_UNSPECIFIED","details":[{"cpeUri":"cpe:o:debian:debian_linux:7","package":"libexempi3","minAffectedVersion":{"epoch":0,"name":"2.5.7","revision":"1","kind":"NORMAL"},"maxAffectedVersion":null,"severityName":"","description":"","fixedLocation":null,"packageType":"","isObsolete":false,"sourceUpdateTime":null}],"cvssV3":null,"windowsDetails":[],"sourceUpdateTime":null}}],"nextPageToken":""}`)
	notesR := objects.ListNotesResponse{}
	json.Unmarshal(data, &notesR)

	tcs := map[string]struct {
		parent        string
		expectedNotes []objects.Note
		returnData    []byte
		expectedErr   string
		statusCode    int
		optsNotes     *objects.ListOpts
		customURL     string
	}{
		"valid input no options": {
			parent:        "grafeastest",
			expectedNotes: notesR.Notes,
			returnData:    data,
			statusCode:    http.StatusOK,
		},
		"valid input with options": {
			parent:        "grafeastest",
			expectedNotes: notesR.Notes,
			returnData:    data,
			statusCode:    http.StatusOK,
			optsNotes: &objects.ListOpts{
				Filter:    optional.NewString("none"),
				PageSize:  optional.NewInt32(10),
				PageToken: optional.NewString("next"),
			},
		},
		"empty parent": {
			parent:        "",
			expectedNotes: notesR.Notes,
			returnData:    data,
			statusCode:    http.StatusOK,
		},
		"json data": {
			parent:        "",
			expectedNotes: notesR.Notes,
			returnData:    data,
			statusCode:    http.StatusOK,
		},
		"invalid json": {
			parent:        "",
			expectedNotes: notesR.Notes,
			returnData:    []byte(`:some { string {`),
			expectedErr:   "invalid character ':' looking for beginning of value",
			statusCode:    http.StatusOK,
		},
		"err status code": {
			parent:        "",
			expectedNotes: notesR.Notes,
			statusCode:    http.StatusInternalServerError,
			expectedErr:   "error getting REST data: " + strconv.Itoa(http.StatusInternalServerError),
		},
	}
	for tc, test := range tcs {
		t.Run(tc, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.statusCode)
				w.Write(test.returnData)
			}))
			defer server.Close()
			if test.customURL == "" {
				test.customURL = server.URL + "/"
			}
			client := NewGrafeasAPIService(test.customURL, "")

			resNotes, err := client.ListNotes(ctx, test.parent, test.optsNotes)

			if err == nil {
				assert.Equal(t, test.expectedNotes, resNotes.Notes)
				assert.Equal(t, nil, err)
			} else {
				assert.Equal(t, test.expectedErr, err.Error())
			}
		})
	}
}

func TestListOccurrences(t *testing.T) {
	ctx := context.Background()
	data := []byte(`{"occurrences":[{"name":"projects/grafeasclienttest/occurrences/0841f19d-16bc-4f83-8a15-640c33fd2e4d","resource":{"name":"","uri":"https://gcr.io/project/image@sha256:foo","contentHash":null},"noteName":"projects/grafeasclienttest/notes/grafeasbuild","kind":"BUILD","remediation":"","createTime":"2020-09-21T20:13:36.120044900Z","updateTime":"2020-09-21T20:13:36.120044900Z","build":{"provenance":{"id":"id","projectId":"ShopifyId","commands":[],"builtArtifacts":[{"checksum":"sha256:71e3e78693c011e59b3fc84940f7672aeeb0a55427b6f5157bd08ab9e9ac746c","id":"some_id1","names":[]}],"createTime":"0001-01-01T00:00:00Z","startTime":"0001-01-01T00:00:00Z","endTime":"0001-01-01T00:00:00Z","creator":"Shopify","logsUri":"some_url","sourceProvenance":{"artifactStorageSourceUri":"","fileHashes":{},"context":{"git":{"url":"https://github.com/Shopify/voucher","revisionId":"q1q"},"labels":{}},"additionalContexts":[]},"triggerId":"","buildOptions":{},"builderVersion":""},"provenanceBytes":""}},{"name":"projects/grafeasclienttest/occurrences/25001d31-f407-419d-83a2-05474ce83e34","resource":{"name":"","uri":"https://gcr.io/project/image@sha256:foo","contentHash":null},"noteName":"projects/grafeasclienttest/notes/grafeasattestation","kind":"ATTESTATION","remediation":"","createTime":"2020-09-21T20:13:36.120087200Z","updateTime":"2020-09-21T20:13:36.120087200Z","attestation":{"attestation":{"pgpSignedAttestation":{"signature":"signature","contentType":"CONTENT_TYPE_UNSPECIFIED","pgpKeyId":"1234"}}}},{"name":"projects/grafeasclienttest/occurrences/332c51c9-4140-4a75-b6b2-b918aba7270d","resource":{"name":"","uri":"https://gcr.io/project/image@sha256:foo","contentHash":null},"noteName":"projects/provider_example/notes/exampleVulnerabilityNote","kind":"VULNERABILITY","remediation":"","createTime":"2020-09-22T19:13:41.283089300Z","updateTime":"2020-09-22T19:13:41.283089300Z","vulnerability":{"type":"","severity":"SEVERITY_UNSPECIFIED","cvssScore":0,"packageIssue":[{"affectedLocation":{"cpeUri":"7","package":"a","version":{"epoch":0,"name":"v1.1.1","revision":"r","kind":"NORMAL"}},"fixedLocation":{"cpeUri":"cpe:/o:debian:debian_linux:7","package":"a","version":{"epoch":0,"name":"namestring","revision":"1","kind":"NORMAL"}},"severityName":""}],"shortDescription":"","longDescription":"","relatedUrls":[],"effectiveSeverity":"SEVERITY_UNSPECIFIED"}},{"name":"projects/grafeasclienttest/occurrences/3690b2f5-cc34-4041-be95-76383cf315cb","resource":{"name":"","uri":"https://gcr.io/project/image@sha256:foo","contentHash":null},"noteName":"projects/grafeasclienttest/notes/grafeasdiscovery","kind":"VULNERABILITY","remediation":"","createTime":"2020-09-21T20:13:36.120127700Z","updateTime":"2020-09-21T20:13:36.120127700Z","vulnerability":{"type":"vulnlow","severity":"SEVERITY_UNSPECIFIED","cvssScore":0,"packageIssue":[{"affectedLocation":{"cpeUri":"uri","package":"package_test","version":{"epoch":0,"name":"v0.1.0","revision":"re","kind":"NORMAL"}},"fixedLocation":null,"severityName":""}],"shortDescription":"","longDescription":"","relatedUrls":[],"effectiveSeverity":"SEVERITY_UNSPECIFIED"}},{"name":"projects/grafeasclienttest/occurrences/7c3c4fc1-8a7e-49df-a399-8771e8778538","resource":{"name":"","uri":"https://gcr.io/project/image@sha256:foo","contentHash":null},"noteName":"projects/grafeasclienttest/notes/grafeasbuild","kind":"BUILD","remediation":"","createTime":"2020-09-21T20:13:36.120065900Z","updateTime":"2020-09-21T20:13:36.120065900Z","build":{"provenance":{"id":"provenceid","projectId":"ShopifyId2","commands":[],"builtArtifacts":[{"checksum":"sha256:71e3e78693c011e59b3fc84940f7672aeeb0a55427b6f5157bd08ab9e9ac746c","id":"some_id2","names":[]}],"createTime":"0001-01-01T00:00:00Z","startTime":"0001-01-01T00:00:00Z","endTime":"0001-01-01T00:00:00Z","creator":"Shopify","logsUri":"some_url2","sourceProvenance":{"artifactStorageSourceUri":"","fileHashes":{},"context":{"git":{"url":"https://github.com/Shopify/voucher","revisionId":"2"},"labels":{}},"additionalContexts":[]},"triggerId":"","buildOptions":{},"builderVersion":""},"provenanceBytes":""}},{"name":"projects/grafeasclienttest/occurrences/8e215cd3-2128-4005-aa54-7ee118620c9a","resource":{"name":"","uri":"https://gcr.io/project/image@sha256:foo","contentHash":null},"noteName":"projects/grafeasclienttest/notes/grafeasvulnerability","kind":"VULNERABILITY","remediation":"","createTime":"2020-09-21T20:13:36.120019600Z","updateTime":"2020-09-21T20:13:36.120019600Z","vulnerability":{"type":"vulnmedium","severity":"MINIMAL","cvssScore":0,"packageIssue":[{"affectedLocation":{"cpeUri":"uri","package":"package_test","version":{"epoch":0,"name":"v0.0.0","revision":"r","kind":"NORMAL"}},"fixedLocation":null,"severityName":""}],"shortDescription":"","longDescription":"","relatedUrls":[],"effectiveSeverity":"SEVERITY_UNSPECIFIED"}},{"name":"projects/grafeasclienttest/occurrences/9cdd69a6-ed0c-4e33-9353-ff4fd83eaf6c","resource":{"name":"","uri":"https://gcr.io/project/image@sha256:foo","contentHash":null},"noteName":"projects/grafeasclienttest/notes/grafeasdiscovery","kind":"DISCOVERY","remediation":"","createTime":"2020-09-21T20:13:36.120148600Z","updateTime":"2020-09-21T20:13:36.120148600Z","discovered":{"discovered":{"continuousAnalysis":"ACTIVE","lastAnalysisTime":null,"analysisStatus":"PENDING","analysisStatusError":null}}},{"name":"projects/grafeasclienttest/occurrences/b1c25d9f-7bc5-4ce5-b0bf-38b31bd03e6e","resource":{"name":"","uri":"https://gcr.io/project/image@sha256:foo","contentHash":null},"noteName":"projects/grafeasclienttest/notes/grafeasdiscovery","kind":"DISCOVERY","remediation":"","createTime":"2020-09-21T20:13:36.120105400Z","updateTime":"2020-09-21T20:13:36.120105400Z","discovered":{"discovered":{"continuousAnalysis":"ACTIVE","lastAnalysisTime":null,"analysisStatus":"FINISHED_SUCCESS","analysisStatusError":null}}}],"nextPageToken":""}`)
	occsR := objects.ListOccurrencesResponse{}
	json.Unmarshal(data, &occsR)

	tcs := map[string]struct {
		parent       string
		expectedOccs []objects.Occurrence
		returnData   []byte
		expectedErr  string
		statusCode   int
		optsNotes    *objects.ListOpts
		customURL    string
	}{
		"valid input": {
			parent:       "grafeastest",
			expectedOccs: occsR.Occurrences,
			returnData:   data,
			statusCode:   http.StatusOK,
		},
		"valid input with options": {
			parent:       "grafeastest",
			expectedOccs: occsR.Occurrences,
			returnData:   data,
			statusCode:   http.StatusOK,
			optsNotes: &objects.ListOpts{
				Filter:    optional.NewString("none"),
				PageSize:  optional.NewInt32(10),
				PageToken: optional.NewString("next"),
			},
		},
		"empty parent": {
			parent:       "",
			expectedOccs: occsR.Occurrences,
			returnData:   data,
			statusCode:   http.StatusOK,
		},
		"json data": {
			parent:       "",
			expectedOccs: occsR.Occurrences,
			returnData:   data,
			statusCode:   http.StatusOK,
		},
		"invalid json": {
			parent:       "",
			expectedOccs: occsR.Occurrences,
			returnData:   []byte(`:{stuff{`),
			statusCode:   http.StatusOK,
			expectedErr:  "invalid character ':' looking for beginning of value",
		},
		"err status code": {
			parent:       "",
			expectedOccs: occsR.Occurrences,
			statusCode:   http.StatusInternalServerError,
			expectedErr:  "error getting REST data: " + strconv.Itoa(http.StatusInternalServerError),
		},
	}
	for tc, test := range tcs {
		t.Run(tc, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(test.statusCode)
				w.Write(test.returnData)
			}))
			defer server.Close()
			if test.customURL == "" {
				test.customURL = server.URL + "/"
			}
			client := NewGrafeasAPIService(test.customURL, "")

			resOccs, err := client.ListOccurrences(ctx, test.parent, test.optsNotes)

			if err == nil {
				assert.Equal(t, test.expectedOccs, resOccs.Occurrences)
				assert.Equal(t, nil, err)
			} else {
				assert.Equal(t, test.expectedErr, err.Error())
			}
		})
	}
}
