# Syscall monkey

Chaos Engineering tool for tampering with syscalls.

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

See [more examples here](./examples).

# Table of contents
- [Syscall monkey](#syscall-monkey)
  - [TL;DR](#tldr)
  - [Quick example](#quick-example)
- [Table of contents](#table-of-contents)
  - [Tutorials](#tutorials)
  - [Installation](#installation)
    - [Binary](#binary)
    - [Use Docker container](#use-docker-container)
    - [Building Docker container](#building-docker-container)
    - [Compatibility](#compatibility)
  - [Advanced usage (SDK)](#advanced-usage-sdk)
  - [TODO](#todo)


## Tutorials

TODO


## Installation

### Binary

You can build the binary from the source:

```sh
git clone https://github.com/seeker89/syscall-monkey.git
cd syscall-monkey
make bin/monkey
./bin/monkey -h
```

### Use Docker container

```sh
docker pull ghcr.io/seeker89/syscall-monkey:v0.0.1rc0
```

```sh
docker run -it ghcr.io/seeker89/syscall-monkey:v0.0.1rc0 -h
Usage:
  monkey [OPTIONS]

Application Options:
  -p, --attach=  Attach to the specified pid
  -t, --target=  Attach to process matching this name
  -c, --config=  Configuration file with desired scenario
  -o, --output=  Write the tracing output to the file (instead of stdout)
  -C, --summary  Show verbose debug information
  -s, --silent   Don't display tracing info

Help Options:
  -h, --help     Show this help message
```

### Building Docker container

```sh
git clone https://github.com/seeker89/syscall-monkey.git
cd syscall-monkey
make build
make run

root@3e14fcd5843c:/# monkey -h
Usage:
  monkey [OPTIONS]

Application Options:
  -p, --attach=  Attach to the specified pid
  -t, --target=  Attach to process matching this name
  -c, --config=  Configuration file with desired scenario
  -o, --output=  Write the tracing output to the file (instead of stdout)
  -C, --summary  Show verbose debug information
  -s, --silent   Don't display tracing info

Help Options:
  -h, --help     Show this help message
```

### Compatibility

Currently, only `Linux` on `x86_64` is supported. If you need arm support, file an issue.


## Advanced usage (SDK)

```go
package main

import (
	"os"

	"github.com/seeker89/syscall-monkey/pkg/config"
	"github.com/seeker89/syscall-monkey/pkg/run"
)

// ExampleManipulator does some random stuff, to illustrate what you can do
type ExampleManipulator struct {
	Count int
}

func (sm *ExampleManipulator) HandleEntry(state SyscallState) SyscallState {
	// change syscall to always be getpid
	state.SyscallCode = 102
	// and also count the entries
	sm.Count++
	return state
}

func (sm *ExampleManipulator) HandleExit(returnValue uint64) uint64 {
	// change the syscall return value on every other call
	if sm.Count%2 == 0 {
		return 0
	}
	return returnValue
}

func main() {
	// parse the config (or hardcode them, if you'd like)
	config := config.ParseCommandLineFlags(os.Args[1:])
	// implement your manipulator
	m := ExampleManipulator{}
	// run the tracer
	run.RunTracer(config, &m)
}
```

## TODO

- [x] auto-generate the mapping of syscall codes to names
- [x] auto-generate the mapping of syscall codes to argument numbers and types
- [x] basic `strace`-like behaviour - start a process, print syscalls and a summary at the end
  - [ ] detach on Ctrl-C
- [x] command line flags handling - drop-in subset for `strace`
- [x] ability to attach to a running process
- [x] add hooks, so that you can implement your own logic and build your own strace
- [x] handle yaml config files
  - [ ] matching and manipulating arguments
- [x] make sure it works in a container as a side car with `SYS_PTRACE`
- [ ] handle tracee's signals
- [ ] ability to find processes by PID, name, or ALL (attach to all processes inside of a container)
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
