package syscall

import (
	"fmt"
	"strings"
	"syscall"
)

func clen(n []byte) int {
	for i := 0; i < len(n); i++ {
		if n[i] == 0 {
			return i
		}
	}
	return len(n)
}

func FormatSyscallExit(regs syscall.PtraceRegs) string {
	ret := int64(regs.Rax)
	if ret >= 0 {
		return fmt.Sprintf(" = %d\n", ret)
	} else {
		ret = -ret
		return fmt.Sprintf(" = -1 (errno %d: %s)\n", ret, syscall.Errno(ret).Error())
	}
}

func FormatSyscallEntry(pid int, regs syscall.PtraceRegs) string {
	code := regs.Orig_rax

	var argsText []string
	for i, argType := range GetSyscallArgumentTypes(code) {
		name := GetSyscallArgumentNames(code)[i]
		value := FormatSyscallArgumentString(pid, argType, ReadSyscallArg(regs, i))
		argsText = append(argsText, fmt.Sprintf("%s=%s", name, value))
	}

	return fmt.Sprintf("%s(%v)", GetSyscallName(code), strings.Join(argsText, ", "))
}

func FormatSyscallArgumentString(pid int, argType string, argVal uint64) string {
	switch argType {
	case "int":
		return fmt.Sprintf("%d", int(argVal))
	case "char __user *", "const char __user *":
		if pid < 1 {
			return "/* char __user * */"
		}
		if argVal == 0 {
			return "NULL"
		}
		var out []byte = make([]byte, 128)
		_, err := syscall.PtracePeekData(pid, uintptr(argVal), out)
		if err != nil {
			panic(err)
		}
		return string(out[:clen(out)])
	case "const struct sigaction __user *", "struct sigaction __user *":
		return fmt.Sprintf("0x%x", argVal)
	default:
		return fmt.Sprintf("%d", argVal)
	}
}
