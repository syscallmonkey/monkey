// +build amd64

package syscall

import "syscall"

func ReadSyscallArg(regs syscall.PtraceRegs, argNo int) uint64 {
	switch argNo {
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
