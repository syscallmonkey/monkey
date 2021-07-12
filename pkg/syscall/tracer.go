package syscall

import (
	"fmt"
	"io"
	"syscall"
)

type Tracer struct {
	Pid     int
	Counter *SyscallCounter
	Out     io.Writer
}

func NewTracer(pid int, out io.Writer) *Tracer {
	t := Tracer{
		Pid:     pid,
		Out:     out,
		Counter: NewSyscallCounter(),
	}
	return &t
}

func (t *Tracer) Loop() {
	var regs syscall.PtraceRegs
	var err error
	var entry bool = true
	// handle the first syscall on its way out - the execve
	syscall.PtraceGetRegs(t.Pid, &regs)
	t.HandleSyscallEntry(regs)
	t.HandleSyscallExit(regs)
	// handle syscalls forever
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
			fmt.Fprintf(t.Out, " = ?\n")
			break
		}
		if entry {
			t.HandleSyscallEntry(regs)
		} else {
			t.HandleSyscallExit(regs)
		}
		entry = !entry
	}
}

func (t *Tracer) HandleSyscallEntry(regs syscall.PtraceRegs) {
	fmt.Fprintf(t.Out, FormatSyscallEntry(t.Pid, regs))
}

func (t *Tracer) HandleSyscallExit(regs syscall.PtraceRegs) {
	t.Counter.Inc(regs.Orig_rax)
	fmt.Fprintf(t.Out, FormatSyscallExit(regs))
}
