package syscall

import (
	"fmt"
	"math/rand"
	"strings"
	"syscall"
	"time"

	config "github.com/syscallmonkey/monkey/pkg/config"
)

type ScenarioManipulator struct {
	Scenario     *config.Scenario
	CurrentRules []*config.SyscallRule
}

func (sm *ScenarioManipulator) HandleEntry(state SyscallState) SyscallState {
	// match the rules and store for the Exit
	sm.CurrentRules = []*config.SyscallRule{}
	for _, rule := range sm.Scenario.Rules {
		if sm.MatchRule(&state, rule) {
			// probability of applying the rule
			rand.Seed(time.Now().UTC().UnixNano())
			if rule.Probability != nil && rand.Float64() > *rule.Probability {
				continue
			}
			sm.CurrentRules = append(sm.CurrentRules, rule)
		}
	}
	// execute the 'before syscall' part of the rules
	for _, rule := range sm.CurrentRules {

		// delay
		if rule.Delay != nil && rule.Delay.Before != nil {
			if duration, err := time.ParseDuration(*rule.Delay.Before); err == nil {
				time.Sleep(duration)
			} else {
				fmt.Printf("\n\nError sleeping '%s'\n\n", err)
			}
		}

		// for blocking, replace the syscall code with something silly
		if rule.Modify != nil && rule.Modify.Block != nil {
			// replace the syscall code with geteuid
			// TODO figure out a better way of doing that
			state.SyscallCode = 107
		}

		if rule.Modify != nil && len(rule.Modify.Args) > 0 {
			for _, modifs := range rule.Modify.Args {
				argTypes := GetSyscallArgumentTypes(state.SyscallCode)
				if modifs.Number < 0 || modifs.Number >= len(argTypes) {
					continue
				}
				if modifs.Int != nil {
					state.Args[modifs.Number] = uint64(*modifs.Int)
				}
				if modifs.Uint != nil {
					state.Args[modifs.Number] = uint64(*modifs.Uint)
				}
				if modifs.String != nil {
					_, err := syscall.PtracePokeData(state.Pid, uintptr(state.Args[modifs.Number]), []byte(*modifs.String))
					if err != nil {
						panic(err)
					}
				}
			}
		}
	}
	return state
}

func (sm *ScenarioManipulator) MatchRule(state *SyscallState, rule *config.SyscallRule) bool {
	match := false

	if rule.Match.Code != nil && *rule.Match.Code == state.SyscallCode {
		match = true
	}
	if rule.Match.Name != "" && strings.Contains(state.SyscallName, rule.Match.Name) {
		match = true
	}

	argsMatch := true
	if match == true && len(rule.Match.Args) > 0 {
		for _, criteria := range rule.Match.Args {
			argTypes := GetSyscallArgumentTypes(state.SyscallCode)
			// ignore faulty ones
			if criteria.Number < 0 || criteria.Number >= len(argTypes) {
				continue
			}
			repr := FormatSyscallArgumentString(state.Pid, argTypes[criteria.Number], state.Args[criteria.Number])
			if criteria.Int != nil && fmt.Sprintf("%d", *criteria.Int) != repr {
				argsMatch = false
				break
			}
			if criteria.Uint != nil && fmt.Sprintf("%d", *criteria.Uint) != repr {
				argsMatch = false
				break
			}
			if criteria.String != nil && *criteria.String != repr {
				argsMatch = false
				break
			}
		}
	}

	return match && argsMatch
}

func (sm *ScenarioManipulator) HandleExit(returnValue uint64) uint64 {
	for _, rule := range sm.CurrentRules {
		if rule.Modify != nil && rule.Modify.Return != nil {
			returnValue = uint64(*rule.Modify.Return)
		}
		if rule.Delay != nil && rule.Delay.After != nil {
			if duration, err := time.ParseDuration(*rule.Delay.After); err == nil {
				time.Sleep(duration)
			} else {
				fmt.Printf("\n\nError sleeping '%s'\n\n", err)
			}
		}
	}
	return returnValue
}
