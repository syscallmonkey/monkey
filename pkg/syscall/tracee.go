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
	pid := cmd.Process.Pid
	var wstatus syscall.WaitStatus
	_, err := syscall.Wait4(pid, &wstatus, 0, nil)
	if err != nil {
		fmt.Printf("Error waiting: %v\n", err)
		panic(err)
	}
	if wstatus.StopSignal() != syscall.SIGTRAP || !wstatus.Stopped() {
		fmt.Printf("Expecting trapped and stopped process\n")
		panic(fmt.Errorf("Expecting trapped and stopped process"))
	}
	fmt.Printf("Started new process pid %d\n", pid)
	return pid, nil
}

// AttachToProcess attaches to an arbitrary process
func AttachToProcess(pid int) error {
	// https://man7.org/linux/man-pages/man2/ptrace.2.html
	// The tracee is sent a SIGSTOP, but
	// will not necessarily have stopped by the completion of
	// this call; use waitpid(2) to wait for the tracee to stop.
	fmt.Printf("Attaching to %d\n", pid)
	err := syscall.PtraceAttach(pid)
	if err != nil {
		fmt.Printf("Error attaching %v\n", err)
		return err
	}
	// wait for it to be stopped
	var wstatus syscall.WaitStatus
	for {
		_, err := syscall.Wait4(pid, &wstatus, 0, nil)
		if err != nil {
			panic(err)
		}
		if wstatus.Stopped() {
			break
		}
	}
	fmt.Printf("Attached to %d\n", pid)
	return nil
}

// DetachProcess detaches from previously attached process
func DetachProcess(pid int) error {
	return syscall.PtraceDetach(pid)
}
