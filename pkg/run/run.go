package run

import (
	"fmt"
	"io"
	"os"

	ps "github.com/mitchellh/go-ps"
	smc "github.com/syscallmonkey/monkey/pkg/config"
	sc "github.com/syscallmonkey/monkey/pkg/syscall"
)

// RunTracer starts a tracer using the provided config and manipulator object
func RunTracer(config *smc.SyscallMonkeyConfig, manipulators []sc.SyscallManipulator) {

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

	if config.AttachPid == 0 && config.TargetName == "" && len(config.TrailingArgs) == 0 {
		fmt.Fprintf(config.OutputFile, "Error: need -p, -t or a command to run\n")
		os.Exit(1)
	}

	if config.AttachPid != 0 {
		err := sc.AttachToProcess(config.AttachPid)
		if err != nil {
			panic(err)
		}
	} else if config.TargetName != "" {
		procs, err := ps.Processes()
		if err != nil {
			panic(err)
		}
		for _, proc := range procs {
			if proc.Executable() == config.TargetName {
				config.AttachPid = proc.Pid()
				break
			}
		}
		if config.AttachPid == 0 {
			panic(fmt.Errorf("No process found for name: %s (%d procs total)", config.TargetName, len(procs)))
		}
		err = sc.AttachToProcess(config.AttachPid)
		if err != nil {
			panic(err)
		}
	} else {
		pid, err := sc.StartTracee(config.TrailingArgs)
		if err != nil {
			panic(err)
		}
		config.AttachPid = pid
	}

	// read the config, if specified
	if config.ConfigPath != "" {
		scenario, err := smc.ParseScenarioFromFile(config.ConfigPath)
		if err != nil {
			panic(err)
		}
		manipulators = append(manipulators, &sc.ScenarioManipulator{
			Scenario: scenario,
		})
	}

	tracer := sc.NewTracer(config.AttachPid, config.OutputFile, manipulators)

	// trace the program until it finishes
	tracer.Loop()

	if config.PrintSummary {
		tracer.Counter.Print(config.OutputFile)
	}
}
