package syscall

type SyscallState struct {
	Arg0        uint64
	Arg1        uint64
	Arg2        uint64
	Arg3        uint64
	Arg4        uint64
	Arg5        uint64
	SyscallCode uint64
	SyscallName string
	Pid         int
}

type SyscallManipulator interface {
	HandleEntry(state SyscallState) SyscallState
	HandleExit(returnValue uint64) uint64
}
