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
	// return value
	case -2:
		return regs.Rax
	// syscall code
	default:
		return regs.Orig_rax
	}
}

func WriteSyscallArg(regs *syscall.PtraceRegs, argNo int, value uint64) {
	switch argNo {
	case 0:
		regs.Rdi = value
	case 1:
		regs.Rsi = value
	case 2:
		regs.Rdx = value
	case 3:
		regs.R10 = value
	case 4:
		regs.R8 = value
	case 5:
		regs.R9 = value
	// return value
	case -2:
		regs.Rax = value
	// syscall code
	default:
		regs.Orig_rax = value
	}
}
