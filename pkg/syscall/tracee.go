package syscall

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// StartTracee starts a new process in Ptrace mode, and awaits the first reutrn
func StartTracee(args []string) (int, error) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Ptrace: true,
	}
	// start and expect error "stop signal: trace/breakpoint trap"
	cmd.Start()
	err := cmd.Wait()
	if err != nil {
		fmt.Printf("Wait returned: %v\n", err)
	}

	return cmd.Process.Pid, nil
}

// AttachToProcess attaches to an arbitrary process
func AttachToProcess(pid int) error {
	return syscall.PtraceAttach(pid)
}
