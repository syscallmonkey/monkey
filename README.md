# syscall-monkey
Chaos Engineering tool for introducing failure into syscalls


## TL;DR

`Syscall Monkey` is like [`strace`](https://man7.org/linux/man-pages/man1/strace.1.html) for fiddling:

- attach and detach processes using [`ptrace`](https://man7.org/linux/man-pages/man2/ptrace.2.html) (Linux only)
- trace their [`syscalls`](https://man7.org/linux/man-pages/man2/syscalls.2.html) - names, arguments, return values
- manipulate `syscalls` (block, change arguments, return value) to simulate failure
- prepare scenarios in a simple `yaml` format
- write custom scenarios using `syscall-monkey` as an SDK

## Quick example

Here's how you can trick `whoami` into changing the user from `root` (0) to `daemon` (1)

```yaml
rules:
  - name: switch geteuid to return daemon
    match:
      name: geteuid
    modify:
      return: 1
```

```sh
root@f34cc94a6b6d:/# whoami
root
root@f34cc94a6b6d:/# monkey -s -c /examples/getuid-user1.yml whoami
daemon
```

## Tutorials

TODO


## Installation

TODO


## TODO

- [x] auto-generate the mapping of syscall codes to names
- [x] auto-generate the mapping of syscall codes to argument numbers and types
- [x] basic `strace`-like behaviour - start a process, print syscalls and a summary at the end
  - [ ] print nicely `umode_t` modes
  - [ ] detach on Ctrl-C
- [x] command line flags handling - drop-in subset for `strace`
- [x] ability to attach to a running process
- [ ] handle tracee's signals
- [ ] ability to find processes by PID, name, or ALL (attach to all processes inside of a container)
- [x] add hooks, so that you can implement your own logic and build your own strace
- [ ] handle yaml config files
- [ ] make sure it works in a container as a side car
- [ ] add an option of HTTP server with an interface to update the stats ()
- [ ] prometheus metrics
- [ ] documentation on how to use
  - [ ] installation
  - [ ] basic strace usage
  - [ ] attaching strace usage
  - [ ] sidecar container for Kubernetes
  - [ ] using the REST api to change the behaviour
