package util

import (
	"flag"
	"os"
	"slices"
	"strconv"
	"strings"
)

var _ = flag.Bool("v", false, "Do verbose logging")

var verbose int = -1

func Verbose() bool {
	return VerboseLevel() > 0
}

func VerboseLevel() int {
	if verbose == -1 {
		verbose = 0
		if slices.Contains(os.Args, "-v") {
			verbose = 1
			return verbose
		}
		value := GetEnv("VERBOSE", "false")
		if strings.ToLower(value) == "true" {
			verbose = 1
		} else if i, err := strconv.Atoi(value); err == nil {
			verbose = i
		}
	}
	return verbose
}
