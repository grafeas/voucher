package main

import (
	"fmt"
	"os"

	flag "github.com/spf13/pflag"
)

func errorf(format string, v ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, format+"\n", v...)
}

func usage() {
	errorf("Usage: %s [flags] <image name>\n", os.Args[0])
	errorf("Flags:")
	flag.PrintDefaults()
}

func main() {
	var server string
	var check string
	var help bool

	flag.StringVarP(&server, "voucher", "v", "localhost", "voucher server address")
	flag.StringVarP(&check, "check", "c", "all", "check to run (\"all\" to run all checks)")
	flag.BoolVarP(&help, "help", "h", false, "print usage message")

	flag.Usage = usage
	flag.Parse()

	if help {
		flag.Usage()
		os.Exit(0)
	}

	if 1 > flag.NArg() {
		errorf("error: not enough arguments")
		flag.Usage()
		os.Exit(1)
	}

	err := lookupAndAttest(server, check, flag.Arg(0))
	if nil != err {
		errorf("error: %s", err)
		os.Exit(1)
	}
}
