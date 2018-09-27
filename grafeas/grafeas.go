package grafeas

import (
	"errors"

	"github.com/docker/distribution/reference"
	"google.golang.org/genproto/googleapis/devtools/containeranalysis/v1beta1/common"
)

var errNoOccurrences = errors.New("no occurrences returned for image")
var errDiscoveriesUnfinished = errors.New("discoveries have not finished processing")

func resourceURL(reference reference.Reference) string {
	return "resourceUrl=\"https://" + reference.String() + "\""
}

func projectPath(project string) string {
	return "projects/" + project
}

func kindFilterStr(reference reference.Reference, kind common.NoteKind) string {
	return resourceURL(reference) + " AND kind=\"" + kind.String() + "\""
}
