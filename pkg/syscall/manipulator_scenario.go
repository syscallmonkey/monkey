package syscall

import (
	config "github.com/seeker89/syscall-monkey/pkg/config"
)

type ScenarioManipulator struct {
	Rules []config.SyscallRule
}

func (sm *ScenarioManipulator) HandleEntry(state SyscallState) SyscallState {
	return state
}

func (sm *ScenarioManipulator) HandleExit(returnValue uint64) uint64 {
	return returnValue
}
