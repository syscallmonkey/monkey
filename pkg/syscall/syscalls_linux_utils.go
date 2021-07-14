package syscall

func GetSyscallName(code uint64) string {
	return codeToName[code]
}

func GetSyscallArgumentTypes(code uint64) []string {
	return codeToArgTypes[code]
}

func GetSyscallArgumentNames(code uint64) []string {
	return codeToArgNames[code]
}
