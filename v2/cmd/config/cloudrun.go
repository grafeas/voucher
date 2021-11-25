package config

import (
	"os"
)

func IsCloudRun() bool {
	return os.Getenv("IS_CLOUDRUN") == "true"
}
