package containeranalysis

import (
	"errors"

	"github.com/docker/distribution/reference"
	grafeas "google.golang.org/genproto/googleapis/grafeas/v1"
)

var errNoOccurrences = errors.New("no occurrences returned for image")
var errDiscoveriesUnfinished = errors.New("discoveries have not finished processing")

func resourceURL(reference reference.Reference) string {
	return "resourceUrl=\"https://" + reference.String() + "\""
}

func projectPath(project string) string {
	return "projects/" + project
}

func kindFilterStr(reference reference.Reference, kind grafeas.NoteKind) string {
	return resourceURL(reference) + " AND kind=\"" + kind.String() + "\""
}
