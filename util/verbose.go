package util

import (
	"flag"
)

var verbose = flag.Bool("v",GetEnvBool("VERBOSE",false),"Do verbose logging")

func Verbose() bool {
	if !flag.Parsed() {
		flag.Parse()
	}

	return *verbose
}