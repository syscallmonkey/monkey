package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type Scenario struct {
	Rules []SyscallRule
}
type SyscallRule struct {
	Name        string             `yaml:"name"`
	Probability float64            `yaml:"probability"`
	Match       SyscallRuleMatch   `yaml:"match"`
	Delay       *SyscallRuleDelay  `yaml:"delay"`
	Modify      *SyscallRuleModify `yaml:"modify"`
}
type SyscallRuleMatch struct {
	Name string            `yaml:"name"`
	Code *uint64           `yaml:"code"`
	Args []SyscallRuleArgs `yaml:"args"`
}
type SyscallRuleDelay struct {
	Before *string `yaml:"before"`
	After  *string `yaml:"after"`
}
type SyscallRuleModify struct {
	Block  *bool             `yaml:"block"`
	Return *int              `yaml:"return"`
	Args   []SyscallRuleArgs `yaml:"args"`
}
type SyscallRuleArgs struct {
	Number int    `yaml:"number"`
	Int    int    `yaml:"int,omitempty"`
	Uint   uint   `yaml:"uint,omitempty"`
	String string `yaml:"string,omitempty"`
}

func ParseScenario(data []byte) (*Scenario, error) {
	s := Scenario{}
	err := yaml.Unmarshal(data, &s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func ParseScenarioFromFile(path string) (*Scenario, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	s, err := ParseScenario(data)
	if err != nil {
		return nil, err
	}
	return s, nil
}
