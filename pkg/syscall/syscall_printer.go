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

func PrintSyscallArgument(argType string, val uint64) string {
	switch argType {
	case "int":
		return fmt.Sprintf("%d", int(val))
	default:
		return fmt.Sprintf("%d", val)
	}
}

func PrintSyscall(regs syscall.PtraceRegs) {
	code := regs.Orig_rax

	var argsText []string
	for i, argType := range GetSyscallArgumentTypes(code) {
		argsText = append(argsText, PrintSyscallArgument(argType, ReadSyscallArg(regs, i)))
	}

	fmt.Printf("%s(%v)", GetSyscallName(code), strings.Join(argsText, ", "))
}
