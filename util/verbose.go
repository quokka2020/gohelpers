package util

import (
	"flag"
)

var verbose *bool = nil

func Verbose() bool {
	if verbose == nil {
		if flag.Parsed() {
			verbose = flag.Bool("v",GetEnvBool("VERBOSE",false),"Do verbose logging")
		} else {
			return false
		}
	}
	return *verbose
}