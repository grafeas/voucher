package grafeas

import (
	"github.com/docker/distribution/reference"
	containeranalysis "google.golang.org/genproto/googleapis/grafeas/v1"
)

func resourceURL(reference reference.Reference) string {
	return "resourceUrl=\"https://" + reference.String() + "\""
}

//ProjectPath defines a project path
func ProjectPath(project string) string {
	return "projects/" + project
}

//KindFilterStr used by containeranalysis
func KindFilterStr(reference reference.Reference, kind containeranalysis.NoteKind) string {
	return resourceURL(reference) + " AND kind=\"" + kind.String() + "\""
}
