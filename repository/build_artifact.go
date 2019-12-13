package repository

// BuildArtifact is a type that describes the artifact info
// related to a build
type BuildArtifact struct {
	ID       string `json:"repository"`
	Checksum string `json:"commit"`
}

func (b *BuildArtifact) String() string {
	str := ""
	if b.ID != "" {
		str += "ID: " + b.ID + "\n"
	}
	if b.Checksum != "" {
		str += "Checksum: " + b.Checksum + "\n"
	}
	return str
}
