package grafeasos

import (
	"errors"

	"github.com/docker/distribution/reference"
	grafeaspb "github.com/grafeas/client-go/0.1.0"
)

var errNoOccurrences = errors.New("no occurrences returned for image")
var errDiscoveriesUnfinished = errors.New("discoveries have not finished processing")

func projectPath(project string) string {
	return "projects/" + project
}

func resourceURL(reference reference.Reference) string {
	return "resourceUrl=\"https://" + reference.String() + "\""
}

func kindFilterStr(reference reference.Reference, kind grafeaspb.V1beta1NoteKind) string {
	return resourceURL(reference) + " AND kind=\"" + string(kind) + "\""
}
