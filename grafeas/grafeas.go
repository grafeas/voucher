package grafeas

import (
	"errors"
	"github.com/Shopify/voucher"
)

var errNoOccurrences = errors.New("no occurrences returned for image")
var errDiscoveriesUnfinished = errors.New("discoveries have not finished processing")

func resourceURL(imageData voucher.ImageData) string {
	return "resourceUrl=\"https://" + imageData.String() + "\""
}

func projectPath(project string) string {
	return "projects/" + project
}

func kindFilterStr(imageData voucher.ImageData, kind voucher.NoteKind) string {
	return resourceURL(imageData) + " AND kind=\"" + kind.String() + "\""
}
