// +build linux

package main

import (
	"fmt"
	"os"

	sc "github.com/seeker89/syscall-monkey/pkg/syscall"
)

var (
	Version, Build string
)

func main() {
	fmt.Printf("Version %s, build %s\n", Version, Build)

	pid, _ := sc.StartTracee(os.Args)
	tracer := sc.NewTracer(pid)
	tracer.Loop()
	tracer.Counter.Print()
}
