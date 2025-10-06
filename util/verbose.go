package util

import (
	"flag"
)

var verbose_flag = flag.Bool("v",false,"Do verbose logging")

var verbose *bool = nil

func Verbose() bool {
	if verbose == nil {
		if *verbose_flag {
			verbose = verbose_flag
		} else if env_file.initialized  {
			*verbose = GetEnvBool("VERBOSE",false)
		} else {
			return false
		}
	}
	return *verbose
}