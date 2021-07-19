package syscall

import (
	"fmt"
	"strings"
	"syscall"
)

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
		printable := make([]byte, 0, 10)
		for i, v := range out {
			if v == 0 {
				break
			}
			switch {
			case v == '\n':
				printable = append(printable, '\\', 'n')
			case v == '\r':
				printable = append(printable, '\\', 'r')
			case v == '\t':
				printable = append(printable, '\\', 't')
			case v >= ' ' && v <= '~':
				printable = append(printable, v)
			default:
				printable = append(printable, []byte(fmt.Sprintf("\\%x", v))...)
			}
			if i > 124 {
				printable = append(printable, []byte("...")...)
				break
			}
		}
		return string(printable)
	case "const struct sigaction __user *", "struct sigaction __user *":
		return fmt.Sprintf("0x%x", argVal)
	default:
		return fmt.Sprintf("%d", argVal)
	}
}
