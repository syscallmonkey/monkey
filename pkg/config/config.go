package config

import (
	"io"
	"os"

	"github.com/jessevdk/go-flags"
)

type SyscallMonkeyConfig struct {
	AttachPid    int    `short:"p" long:"attach" description:"Attach to the specified pid"`
	TargetName   string `short:"t" long:"target" description:"Attach to process matching this name"`
	ConfigPath   string `short:"c" long:"config" description:"Configuration file with desired scenario"`
	OutputPath   string `short:"o" long:"output" description:"Write the tracing output to the file (instead of stdout)"`
	PrintSummary bool   `short:"C" long:"summary" description:"Show verbose debug information"`
	Silent       bool   `short:"s" long:"silent" description:"Don't display tracing info"`
	TrailingArgs []string
	OutputFile   io.Writer
	Version      string
	Build        string
	TraceeStdin  io.Reader
	TraceeStdout io.Writer
	TraceeStderr io.Writer
	Pid          int
}

func ParseCommandLineFlags(args []string) *SyscallMonkeyConfig {
	config := SyscallMonkeyConfig{}
	trailing, err := flags.ParseArgs(&config, args)
	config.TrailingArgs = trailing
	if err != nil {
		os.Exit(1)
	}
	return &config
}
