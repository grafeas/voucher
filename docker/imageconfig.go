package docker

import (
	dockerTypes "github.com/docker/docker/api/types"
)

// ImageConfig is a structure that wraps a config manifest.
type ImageConfig struct {
	ContainerConfig dockerTypes.ExecConfig `json:"container_config"`
}

// RunsAsRoot returns true if the image the ImageConfig is associated with
// runs as root (user 0).
func (config *ImageConfig) RunsAsRoot() bool {
	user := config.ContainerConfig.User

	return ("" == user || "root" == user || "0:0" == user || "0" == user)
}
