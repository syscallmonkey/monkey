// +build linux

package main

import (
	"fmt"
	"os"
	"syscall"

	sc "github.com/seeker89/syscall-monkey/pkg/syscall"
)

var (
	Version, Build string
)

func main() {
	fmt.Printf("Version %s, build %s\n", Version, Build)

	pid, _ := sc.StartTracee(os.Args)
	incall := false
	counter := sc.NewSyscallCounter()

	var regs syscall.PtraceRegs

	for {
		err := syscall.PtraceGetRegs(pid, &regs)
		if err != nil {
			fmt.Printf(" = ?\n")
			break
		}
		code := regs.Orig_rax
		name := sc.GetSyscallName(code)
		if incall {
			fmt.Printf("%s()", name)
			counter.Inc(code)
		} else {
			fmt.Printf(" = %d\n", regs.Rax)
		}
		incall = !incall

		err = syscall.PtraceSyscall(pid, 0)
		if err != nil {
			panic(err)
		}

		_, err = syscall.Wait4(pid, nil, 0, nil)
		if err != nil {
			panic(err)
		}

	}
	counter.Print()
}
