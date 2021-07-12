// +build linux

package main

import (
	"fmt"
	"os"

	smc "github.com/seeker89/syscall-monkey/pkg/config"
	sc "github.com/seeker89/syscall-monkey/pkg/syscall"
)

var (
	Version, Build string
)

func main() {

	// parse the flags
	config := smc.ParseCommandLineFlags(os.Args[1:])

	// figure out where to direct the output
	config.OutputFile = os.Stdout
	if config.OutputPath != "" {
		f, err := os.OpenFile(config.OutputPath, os.O_CREATE|os.O_WRONLY, 0660)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		config.OutputFile = f
	}

	// just a reminder
	fmt.Fprintf(config.OutputFile, "Version %s, build %s\n", Version, Build)

	if config.AttachPid == 0 {
		if len(config.TrailingArgs) == 0 {
			fmt.Fprintf(config.OutputFile, "Error: need either -p or a command to run\n")
			os.Exit(1)
		}
		pid, err := sc.StartTracee(config.TrailingArgs)
		if err != nil {
			panic(err)
		}
		config.AttachPid = pid
	} else {
		err := sc.AttachToProcess(config.AttachPid)
		if err != nil {
			panic(err)
		}
	}

	tracer := sc.NewTracer(config.AttachPid, config.OutputFile)

	// trace the program until it finishes
	tracer.Loop()

	if config.PrintSummary {
		tracer.Counter.Print(config.OutputFile)
	}
}
