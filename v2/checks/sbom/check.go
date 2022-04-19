package sbom

import (
	"github.com/spdx/tools-golang"
)

// There are two supported versions in
// https://github.com/spdx/tools-golang/tree/main/spdx
// i.e. 2.1 and 2.2; we'll need to support both for spdx

type SBOMResponse struct {
	SPDXVersion string `json:"spdxversion"`
}

type SBOM interface {
	GetSBOM() string
	GetVersion() string
}

func (s *SBOMResponse) GetVersion() string {
	return s.SPDXVersion
}

func (s *SBOMResponse) GetSBOM() string {

}

// CreationInfo2_1
// https://github.com/spdx/tools-golang/blob/8e09d22f514a2eeee9a3044fb47d0cb6dbc74a6a/spdx/creation_info.go
