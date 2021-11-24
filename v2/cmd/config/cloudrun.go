package config

import (
	"os"
	"strconv"
)

func IsCloudRun() bool {
	isCloudRun := os.Getenv("IS_CLOUDRUN")
	isCloudRunBool, err := strconv.ParseBool(isCloudRun)
	if err != nil {
		return false
	}
	return isCloudRunBool
}
