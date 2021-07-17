// +build linux

package main

import (
	"os"

	smc "github.com/syscallmonkey/monkey/pkg/config"
	smr "github.com/syscallmonkey/monkey/pkg/run"
)

var (
	Version, Build string
)

func main() {

	// parse the flags
	config := smc.ParseCommandLineFlags(os.Args[1:])

	// provide the build metadata for nice output
	config.Version = Version
	config.Build = Build

	// run the tracer
	smr.RunTracer(config, nil)
}
