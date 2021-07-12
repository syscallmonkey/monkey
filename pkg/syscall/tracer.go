package syscall

import (
	"fmt"
	"io"
	"syscall"
)

type Tracer struct {
	Pid         int
	Counter     *SyscallCounter
	Out         io.Writer
	Manipulator SyscallManipulator
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
	if t.Manipulator != nil {
		code := ReadSyscallArg(regs, -1)
		state := SyscallState{
			Arg0:        ReadSyscallArg(regs, 0),
			Arg1:        ReadSyscallArg(regs, 1),
			Arg2:        ReadSyscallArg(regs, 2),
			Arg3:        ReadSyscallArg(regs, 3),
			Arg4:        ReadSyscallArg(regs, 4),
			Arg5:        ReadSyscallArg(regs, 5),
			SyscallCode: code,
			SyscallName: GetSyscallName(code),
			Pid:         t.Pid,
		}
		newState := t.Manipulator.HandleEntry(state)
		if newState != state {
			newRegs := regs
			WriteSyscallArg(&newRegs, -1, newState.SyscallCode)
			WriteSyscallArg(&newRegs, 0, newState.Arg0)
			WriteSyscallArg(&newRegs, 1, newState.Arg1)
			WriteSyscallArg(&newRegs, 2, newState.Arg2)
			WriteSyscallArg(&newRegs, 3, newState.Arg3)
			WriteSyscallArg(&newRegs, 4, newState.Arg4)
			WriteSyscallArg(&newRegs, 5, newState.Arg5)
			syscall.PtraceSetRegs(t.Pid, &newRegs)
		}
	}
}

func (t *Tracer) HandleSyscallExit(regs syscall.PtraceRegs) {
	t.Counter.Inc(regs.Orig_rax)
	fmt.Fprintf(t.Out, FormatSyscallExit(regs))
	if t.Manipulator != nil {
		ret := ReadSyscallArg(regs, -2)
		newRet := t.Manipulator.HandleExit(ret)
		if newRet != ret {
			newRegs := regs
			WriteSyscallArg(&newRegs, -2, newRet)
			syscall.PtraceSetRegs(t.Pid, &newRegs)
		}
	}
}
