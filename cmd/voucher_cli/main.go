package main

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	rootCmd.Version = version
	rootCmd.Execute()
}
