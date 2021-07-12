package syscall

import (
	"fmt"
	"strings"
	"syscall"
)

func ReadSyscallArg(regs syscall.PtraceRegs, arg int) uint64 {
	switch arg {
	case 0:
		return regs.Rdi
	case 1:
		return regs.Rsi
	case 2:
		return regs.Rdx
	case 3:
		return regs.R10
	case 4:
		return regs.R8
	case 5:
		return regs.R9
	default:
		return regs.Orig_rax
	}
}

func clen(n []byte) int {
	for i := 0; i < len(n); i++ {
		if n[i] == 0 {
			return i
		}
	}
	return len(n)
}

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
	t.PrintSyscall(regs)
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
			t.PrintSyscall(regs)
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

func (t *Tracer) PrintSyscall(regs syscall.PtraceRegs) {
	code := regs.Orig_rax

	var argsText []string
	for i, argType := range GetSyscallArgumentTypes(code) {
		argsText = append(argsText, t.PrintSyscallArgument(argType, ReadSyscallArg(regs, i)))
	}

	fmt.Printf("%s(%v)", GetSyscallName(code), strings.Join(argsText, ", "))
}

func (t *Tracer) PrintSyscallArgument(argType string, val uint64) string {
	switch argType {
	case "int":
		return fmt.Sprintf("%d", int(val))
	case "const char __user *":
		if val == 0 {
			return "NULL"
		}
		var out []byte = make([]byte, 128)
		_, err := syscall.PtracePeekData(t.Pid, uintptr(val), out)
		if err != nil {
			panic(err)
		}
		return string(out[:clen(out)])
	default:
		return fmt.Sprintf("%d", val)
	}
}
