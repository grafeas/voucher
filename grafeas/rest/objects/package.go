package objects

//VersionKind https://github.com/grafeas/client-go/blob/master/0.1.0/model_version_version_kind.go
type VersionKind string

//PackageArchitecture https://github.com/grafeas/client-go/blob/master/0.1.0/model_package_architecture.go
type PackageArchitecture string

//consts
const (
	VersionKindUnspecified VersionKind = "VERSION_KIND_UNSPECIFIED"
	VersionKindNormal      VersionKind = "NORMAL"
	VersionKindMinimum     VersionKind = "MINIMUM"
	VVersionKindMaximum    VersionKind = "MAXIMUM"

	PackageArchitectureUnspecified PackageArchitecture = "ARCHITECTURE_UNSPECIFIED"
	PackageArchitectureX86         PackageArchitecture = "X86"
	PackageArchitectureX64         PackageArchitecture = "X64"
)

//Package https://github.com/grafeas/client-go/blob/master/0.1.0/model_package_package.go
type Package struct {
	Name         string                `json:"name,omitempty"` //required
	Distribution []PackageDistribution `json:"distribution,omitempty"`
}

//PackageDistribution https://github.com/grafeas/client-go/blob/master/0.1.0/model_package_distribution.go
type PackageDistribution struct {
	CpeURI        string               `json:"cpeUri,omitempty"` //required
	Architecture  *PackageArchitecture `json:"architecture,omitempty"`
	LatestVersion *PackageVersion      `json:"latestVersion,omitempty"`
	Maintainer    string               `json:"maintainer,omitempty"`
	URL           string               `json:"url,omitempty"`
	Description   string               `json:"description,omitempty"`
}

//PackageVersion https://github.com/grafeas/client-go/blob/master/0.1.0/model_package_version.go
type PackageVersion struct {
	Epoch    int32        `json:"epoch,omitempty"`
	Name     string       `json:"name,omitempty"` //required only when version kind is NORMAL
	Revision string       `json:"revision,omitempty"`
	Kind     *VersionKind `json:"kind,omitempty"` //required
}

//PackageDetails https://github.com/grafeas/client-go/blob/master/0.1.0/model_v1beta1package_details.go
type PackageDetails struct {
	Installation *PackageInstallation `json:"installation,omitempty"`
}

//PackageInstallation https://github.com/grafeas/client-go/blob/master/0.1.0/model_package_installation.go
type PackageInstallation struct {
	Name     string            `json:"name,omitempty"`     //output only
	Location []PackageLocation `json:"location,omitempty"` //required
}

//PackageLocation https://github.com/grafeas/client-go/blob/master/0.1.0/model_v1beta1package_location.go
type PackageLocation struct {
	CpeURI  string          `json:"cpeUri,omitempty"` //required
	Version *PackageVersion `json:"version,omitempty"`
	Path    string          `json:"path,omitempty"`
}
