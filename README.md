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
- [x] ability to find processes by PID or name
- [x] automate publishing docker image
- [ ] figure out why the arguments for the first syscall don't work
- [x] change how the strings and buffers are printed - handle `\n`, `\r`, `\t` and non-printable characters
- [ ] multi-process support
  - [ ] ability to detect and follow children
  - [ ] ability to find multiple processes
- [ ] HTTP server
  - [ ] update the config on the fly
  - [ ] prometheus metrics
- [ ] write unit test coverage LOL
- [ ] documentation on how to use
  - [x] installation
  - [x] basic strace-like usage
  - [x] attaching strace usage
  - [ ] plenty of cool examples
  - [ ] sidecar container for Kubernetes
  - [ ] using the REST api to change the behaviour
