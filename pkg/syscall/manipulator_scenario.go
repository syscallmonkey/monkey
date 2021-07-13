package syscall

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	config "github.com/seeker89/syscall-monkey/pkg/config"
)

type ScenarioManipulator struct {
	Scenario     *config.Scenario
	CurrentRules []*config.SyscallRule
}

func (sm *ScenarioManipulator) HandleEntry(state SyscallState) SyscallState {
	// match the rules and store for the Exit
	sm.CurrentRules = []*config.SyscallRule{}
	for _, rule := range sm.Scenario.Rules {
		if sm.MatchRule(&state, &rule) {
			// probability of applying the rule
			if rule.Probability != nil && rand.Float64() > *rule.Probability {
				continue
			}
			sm.CurrentRules = append(sm.CurrentRules, &rule)
		}
	}
	// execute the 'before syscall' part of the rules
	for _, rule := range sm.CurrentRules {

		fmt.Printf("\n\nExecuting rule '%s'\n\n", rule.Name)

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

	// TODO handle argument matching

	return match
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
