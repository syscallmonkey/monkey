package syscall

import (
	"fmt"
	"syscall"
)

type Tracer struct {
	Pid     int
	Counter *SyscallCounter
}

func NewTracer(pid int) *Tracer {
	t := Tracer{
		Pid:     pid,
		Counter: NewSyscallCounter(),
	}
	return &t
}

func (t *Tracer) Loop() {
	var regs syscall.PtraceRegs
	var err error
	var code uint64
	var incall bool
	// handle the first syscall on its way out - the execve
	syscall.PtraceGetRegs(t.Pid, &regs)
	code = regs.Orig_rax
	t.Counter.Inc(code)
	fmt.Printf("%s() = %d", GetSyscallName(code), regs.Rax)
	incall = true
	// handle all the other syscalls
	for {
		err = syscall.PtraceSyscall(t.Pid, 0)
		if err != nil {
			panic(err)
		}
		_, err = syscall.Wait4(t.Pid, nil, 0, nil)
		if err != nil {
			panic(err)
		}
		err = syscall.PtraceGetRegs(t.Pid, &regs)
		if err != nil {
			fmt.Printf(" = ?\n")
			break
		}
		if incall {
			code = regs.Orig_rax
			fmt.Printf("%s()", GetSyscallName(code))
			t.Counter.Inc(code)
		} else {
			fmt.Printf(" = %d\n", regs.Rax)
		}
		incall = !incall
	}
}
