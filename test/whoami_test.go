package test

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	smc "github.com/syscallmonkey/monkey/pkg/config"
	smrun "github.com/syscallmonkey/monkey/pkg/run"
)

//go:embed examples/getuid-user1.yml
var getUidUser1 string

func GetUserNameById(id int) (string, error) {
	// return user from etc/passwd
	content, err := ioutil.ReadFile("/etc/passwd")
	if err != nil {
		return "", fmt.Errorf("error reading /etc/passwd: %v", err)
	}
	var name string
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		items := strings.Split(line, ":")
		if len(items) == 7 && items[2] == "1" {
			name = items[0]
		}
	}
	return name, nil
}

func TestChangeGetUidReturnValue(t *testing.T) {

	expectedUser, err := GetUserNameById(1)
	if err != nil {
		t.Errorf("Error getting user name %v", err)
	}

	// run the thing
	output := strings.Builder{}
	config := smc.SyscallMonkeyConfig{
		ConfigPath:   "./examples/getuid-user1.yml",
		TrailingArgs: []string{"whoami"},
		Silent:       true,
		TraceeStdout: &output,
	}
	smrun.RunTracer(&config, nil)

	outputTrimmed := strings.TrimSpace(output.String())

	if outputTrimmed != expectedUser {
		t.Errorf("Expected '%s', got '%s' (%s)", expectedUser, outputTrimmed, output.String())
	}
}
