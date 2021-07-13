package run

import (
	"fmt"
	"io"
	"os"

	smc "github.com/seeker89/syscall-monkey/pkg/config"
	sc "github.com/seeker89/syscall-monkey/pkg/syscall"
)

// RunTracer starts a tracer using the provided config and manipulator object
func RunTracer(config *smc.SyscallMonkeyConfig, manipulator sc.SyscallManipulator) {

	// figure out where to direct the output
	if config.Silent {
		config.OutputFile = io.Discard
	} else if config.OutputFile == nil && config.OutputPath != "" {
		f, err := os.OpenFile(config.OutputPath, os.O_CREATE|os.O_WRONLY, 0660)
		if err != nil {
			panic(err)
		}
		defer f.Close()
		config.OutputFile = f
	} else {
		config.OutputFile = os.Stdout
	}

	// just a reminder
	fmt.Fprintf(config.OutputFile, "Version %s, build %s\n", config.Version, config.Build)

	if config.AttachPid == 0 {
		if len(config.TrailingArgs) == 0 {
			fmt.Fprintf(config.OutputFile, "Error: need either -p or a command to run\n")
			os.Exit(1)
		}
		pid, err := sc.StartTracee(config.TrailingArgs)
		if err != nil {
			panic(err)
		}
		config.AttachPid = pid
	} else {
		err := sc.AttachToProcess(config.AttachPid)
		if err != nil {
			panic(err)
		}
	}

	// read the config, if specified
	if config.ConfigPath != "" {
		scenario, err := smc.ParseScenarioFromFile(config.ConfigPath)
		if err != nil {
			panic(err)
		}
		manipulator = &sc.ScenarioManipulator{
			Scenario: scenario,
		}
	}

	tracer := sc.NewTracer(config.AttachPid, config.OutputFile, manipulator)

	// trace the program until it finishes
	tracer.Loop()

	if config.PrintSummary {
		tracer.Counter.Print(config.OutputFile)
	}
}
