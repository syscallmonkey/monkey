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

func NewTracer(pid int, out io.Writer, manipulator SyscallManipulator) *Tracer {
	t := Tracer{
		Pid:         pid,
		Out:         out,
		Manipulator: manipulator,
		Counter:     NewSyscallCounter(),
	}
	return &t
}

func (t *Tracer) Loop() {
	var regs syscall.PtraceRegs
	var err error
	var entry bool = true
	// Set options to detect our syscalls
	// https://man7.org/linux/man-pages/man2/ptrace.2.html
	// PTRACE_O_TRACESYSGOOD (since Linux 2.4.6)
	// When delivering system call traps, set bit 7 in the
	// signal number (i.e., deliver SIGTRAP|0x80).  This
	// makes it easy for the tracer to distinguish normal
	// traps from those caused by a system call.
	err = syscall.PtraceSetOptions(t.Pid, syscall.PTRACE_O_TRACESYSGOOD)
	if err != nil {
		fmt.Printf("Error setting options %v\n", err)
		panic(err)
	}
	// TODO FIXME the call parameters are empty
	// handle the first syscall on its way out - the execve
	err = syscall.PtraceGetRegs(t.Pid, &regs)
	if err != nil {
		panic(err)
	}
	t.HandleSyscallEntry(regs)
	t.HandleSyscallExit(regs)
	// handle syscalls forever
	for {
		for {
			err = syscall.PtraceSyscall(t.Pid, 0)
			if err != nil {
				fmt.Printf("error: %v, pid: %d\n", err, t.Pid)
				panic(err)
			}
			var wstatus syscall.WaitStatus
			_, err := syscall.Wait4(t.Pid, &wstatus, 0, nil)
			if err != nil {
				panic(err)
			}
			// on execve
			if wstatus.StopSignal() == syscall.SIGTRAP {
				fmt.Printf("GOT a SIGTRAP (execve)\n")
			}
			// stopped and stopped for us (syscall.PTRACE_O_TRACESYSGOOD)
			if wstatus.Stopped() && wstatus.StopSignal()&0x80 == 0x80 {
				break
			}
			// or terminated
			if wstatus.Exited() {
				fmt.Fprintf(t.Out, "\n--- program exited ---\n")
				return
			}
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
			SyscallCode: code,
			SyscallName: GetSyscallName(code),
			Pid:         t.Pid,
		}
		for i := range state.Args {
			state.Args[i] = ReadSyscallArg(regs, i)
		}
		newState := t.Manipulator.HandleEntry(state)
		if newState != state {
			newRegs := regs
			WriteSyscallArg(&newRegs, -1, newState.SyscallCode)
			for i := range newState.Args {
				WriteSyscallArg(&newRegs, i, newState.Args[i])
			}
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
