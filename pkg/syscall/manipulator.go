package syscall

type SyscallState struct {
	Args        [6]uint64
	SyscallCode uint64
	SyscallName string
	Pid         int
}

type SyscallManipulator interface {
	HandleEntry(state SyscallState) SyscallState
	HandleExit(returnValue uint64) uint64
}
