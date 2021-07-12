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
	fmt.Printf("Version %s, build %s\n", Version, Build)

	config := smc.ParseCommandLineFlags(os.Args[1:])

	if config.AttachPid == 0 {
		pid, err := sc.StartTracee(os.Args)
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

	tracer := sc.NewTracer(config.AttachPid)
	tracer.Loop()
	tracer.Counter.Print()
}
