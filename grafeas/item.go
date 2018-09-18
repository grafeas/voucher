package grafeas

import (
	"strings"

	containeranalysispb "google.golang.org/genproto/googleapis/devtools/containeranalysis/v1alpha1"
)

// Item implements a MetadataItem.
type Item struct {
	Occurrence *containeranalysispb.Occurrence // The containeranalysispb Occurrence this Item wraps.
}

// Name returns the name of the group of Item.
func (item *Item) Name() string {
	return item.Occurrence.NoteName
}

// Kind returns the kind of item.
func (item *Item) Kind() string {
	noteName := strings.Split(item.Occurrence.NoteName, "/")
	return noteName[len(noteName)-1]
}

// String returns a string version of this Item.
func (item *Item) String() string {
	return item.Occurrence.String()
}
