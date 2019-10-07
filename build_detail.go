package voucher

import (
	"strings"
)

// BuildDetail is a type that describes the details/metadata info
// related to a build
type BuildDetail struct {
	RepositoryURL string          `json:"repository"`
	Commit        string          `json:"commit"`
	BuildCreator  string          `json:"build_creator"`
	BuildURL      string          `json:"build_url"`
	ProjectID     string          `json:"project_id"`
	Artifacts     []BuildArtifact `json:"artifacts"`
}

func (b *BuildDetail) String() string {
	str := ""
	if b.RepositoryURL != "" {
		str += "RepositoryURL: " + b.RepositoryURL + "\n"
	}
	if b.Commit != "" {
		str += "Commit: " + b.Commit + "\n"
	}
	if b.BuildCreator != "" {
		str += "BuildCreator: " + b.BuildCreator + "\n"
	}
	if b.BuildURL != "" {
		str += "BuildURL: " + b.BuildURL + "\n"
	}
	if b.ProjectID != "" {
		str += "ProjectID: " + b.ProjectID + "\n"
	}
	str_artifacts := ""
	for _, val := range b.Artifacts {
		if val.String() != "" {
			str_artifacts = strings.Join([]string{str_artifacts, val.String()}, ", ")
		}
	}
	if str_artifacts != "" {
		str += "Artifacts: " + str_artifacts + "\n"
	}
	return str
}
