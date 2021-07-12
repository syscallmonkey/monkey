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
	var incall bool
	// handle the first syscall on its way out - the execve
	syscall.PtraceGetRegs(t.Pid, &regs)
	t.Counter.Inc(regs.Orig_rax)
	PrintSyscall(regs)
	fmt.Printf(" = %d\n", regs.Rax)
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
			PrintSyscall(regs)
			t.Counter.Inc(regs.Orig_rax)
		} else {
			ret := int64(regs.Rax)
			if ret >= 0 {
				fmt.Printf(" = %d\n", ret)
			} else {
				ret = -ret
				fmt.Printf(" = -1 (errno %d: %s)\n", ret, syscall.Errno(ret).Error())
			}
		}
		incall = !incall
	}
}
