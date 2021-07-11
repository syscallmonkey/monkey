package syscall

import (
	"fmt"
	"syscall"
)

type Tracer struct {
	Incall  bool
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
	for {
		var regs syscall.PtraceRegs
		err := syscall.PtraceGetRegs(t.Pid, &regs)
		if err != nil {
			fmt.Printf(" = ?\n")
			break
		}
		code := regs.Orig_rax
		name := GetSyscallName(code)
		if t.Incall {
			fmt.Printf("%s()", name)
			t.Counter.Inc(code)
		} else {
			fmt.Printf(" = %d\n", regs.Rax)
		}
		t.Incall = !t.Incall

		err = syscall.PtraceSyscall(t.Pid, 0)
		if err != nil {
			panic(err)
		}
		_, err = syscall.Wait4(t.Pid, nil, 0, nil)
		if err != nil {
			panic(err)
		}
	}
}
