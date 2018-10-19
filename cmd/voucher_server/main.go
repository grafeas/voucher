package main

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	serverCmd.Version = version
	serverCmd.Execute()
}
