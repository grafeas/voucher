package objects

// VersionKind based on
// https://github.com/grafeas/client-go/blob/master/0.1.0/model_version_version_kind.go
type VersionKind string

// consts
const (
	VersionKindUnspecified VersionKind = "VERSION_KIND_UNSPECIFIED"
	VersionKindNormal      VersionKind = "NORMAL"
	VersionKindMinimum     VersionKind = "MINIMUM"
	VVersionKindMaximum    VersionKind = "MAXIMUM"
)

// Package based on
// https://github.com/grafeas/client-go/blob/master/0.1.0/model_package_package.go
type Package struct {
	Name string `json:"name,omitempty"` //required
}

// PackageVersion based on
// https://github.com/grafeas/client-go/blob/master/0.1.0/model_package_version.go
type PackageVersion struct {
	Epoch    int32        `json:"epoch,omitempty"`
	Name     string       `json:"name,omitempty"` //required only when version kind is NORMAL
	Revision string       `json:"revision,omitempty"`
	Kind     *VersionKind `json:"kind,omitempty"` //required
}
