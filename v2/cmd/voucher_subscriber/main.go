package main

import (
	"fmt"
	"os"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	subscriberCmd.Version = version

	if err := subscriberCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
