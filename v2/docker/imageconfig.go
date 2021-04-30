package docker

import (
	dockerTypes "github.com/docker/docker/api/types"
)

// ImageConfig represents an Docker image configuration. This presently just
// allows us to verify if an image runs as root or not.
type ImageConfig interface {
	// RunsAsRoot returns true if the passed image will run as the root user.
	RunsAsRoot() bool
}

type imageConfig struct {
	dockerTypes.ExecConfig
}

// RunsAsRoot returns true if the image will run as the root user.
func (config *imageConfig) RunsAsRoot() bool {
	user := config.User

	return ("" == user || "root" == user || "0:0" == user || "0" == user)
}
