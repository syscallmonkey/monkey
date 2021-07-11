// +build linux

package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	sc "github.com/seeker89/syscall-monkey/pkg/syscall"
)

var (
	Version, Build string
)

func main() {
	fmt.Printf("Version %s, build %s\n", Version, Build)

	fmt.Printf("Run %v\n", os.Args[1:])

	cmd := exec.Command(os.Args[1], os.Args[2:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Ptrace: true,
	}

	cmd.Start()
	err := cmd.Wait()
	if err != nil {
		fmt.Printf("Wait returned: %v\n", err)
	}

	pid := cmd.Process.Pid
	incall := false
	counter := sc.NewSyscallCounter()

	var regs syscall.PtraceRegs

	for {
		err = syscall.PtraceGetRegs(pid, &regs)
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
