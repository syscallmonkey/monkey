# Syscall monkey

Chaos Engineering tool for tampering with syscalls.

To get started, check out the docs: https://syscallmonkey.github.io/


## TODO

- [x] auto-generate the mapping of syscall codes to names
- [x] auto-generate the mapping of syscall codes to argument numbers and types
- [x] basic `strace`-like behaviour - start a process, print syscalls and a summary at the end
  - [ ] clean detach on Ctrl-C
- [x] command line flags handling - drop-in subset for `strace`
- [x] ability to attach to a running process
- [x] add hooks, so that you can implement your own logic and build your own strace
- [x] handle yaml config files
  - [x] matching and manipulating arguments
- [x] make sure it works in a container as a side car with `SYS_PTRACE`
- [x] handle tracee's signals
- [ ] ability to find processes by PID, name, or ALL (attach to all processes inside of a container)
- [ ] ability to detect and follow children
- [ ] HTTP server
  - [ ] update the config on the fly
  - [ ] prometheus metrics
- [ ] write unit test coverage LOL
- [ ] publish an image to Docker Hub
- [ ] documentation on how to use
  - [ ] installation
  - [ ] basic strace-like usage
  - [ ] plenty of cool examples
  - [ ] attaching strace usage
  - [ ] sidecar container for Kubernetes
  - [ ] using the REST api to change the behaviour
