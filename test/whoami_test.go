package test

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	smc "github.com/syscallmonkey/monkey/pkg/config"
	smrun "github.com/syscallmonkey/monkey/pkg/run"
)

var EXAMPLE_GETUID_USER1_PATH string = "./examples/getuid-user1.yml"

func GetUserNameById(id int) (string, error) {
	// return user from etc/passwd
	content, err := ioutil.ReadFile("/etc/passwd")
	if err != nil {
		return "", fmt.Errorf("error reading /etc/passwd: %v", err)
	}
	var name string
	idString := fmt.Sprintf("%d", id)
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		items := strings.Split(line, ":")
		if len(items) == 7 && items[2] == idString {
			name = items[0]
		}
	}
	return name, nil
}

func TestChangeGetUidReturnValue(t *testing.T) {

	// check the that running as regular user matches the output
	regularUser, err := GetUserNameById(0)
	if err != nil {
		t.Errorf("Error getting user name %v", err)
	}
	output := strings.Builder{}
	config := smc.SyscallMonkeyConfig{
		TrailingArgs: []string{"whoami"},
		Silent:       true,
		TraceeStdout: &output,
	}
	smrun.RunTracer(&config, nil)

	regularOutput := strings.TrimSpace(output.String())

	if regularOutput != regularUser {
		t.Errorf("Expected '%s', got '%s' (%s)", regularUser, regularOutput, output.String())
	}

	// get user 1
	modifiedUser, err := GetUserNameById(1)
	if err != nil {
		t.Errorf("Error getting user name %v", err)
	}

	if regularUser == modifiedUser {
		t.Errorf("Expected '%s' to be different from '%s'", regularUser, modifiedUser)
	}

	// check that running with the scenario it returns a different user
	output.Reset()
	config = smc.SyscallMonkeyConfig{
		ConfigPath:   EXAMPLE_GETUID_USER1_PATH,
		TrailingArgs: []string{"whoami"},
		Silent:       true,
		TraceeStdout: &output,
	}
	smrun.RunTracer(&config, nil)

	modifiedOutput := strings.TrimSpace(output.String())

	if modifiedOutput != modifiedUser {
		t.Errorf("Expected '%s', got '%s' (%s)", modifiedUser, modifiedOutput, output.String())
	}
}

func TestChangeEtcPasswd(t *testing.T) {

	template := "%s:x:0:0:root:/root:/bin/bash\n"

	var tests = []struct {
		scenario, path string
	}{
		{"openat-etc-passwd.yml", "/tmp/passwd"},
		{"openat-etc-passwd-short.yml", "/tmp/p"},
		{"openat-etc-passwd-long.yml", "/tmp/pretty-long-passwd"},
	}
	for i, tt := range tests {
		testname := fmt.Sprintf("Scenario %s (%s)", tt.scenario, tt.path)
		t.Run(testname, func(t *testing.T) {

			expectedUser := fmt.Sprintf("user-%d", i)
			content := fmt.Sprintf(template, expectedUser)

			// prep the file
			err := ioutil.WriteFile(tt.path, []byte(content), 0644)
			if err != nil {
				t.Errorf("Error writing temp file %s: %v", tt.path, err)
			}
			defer func() {
				err = os.Remove(tt.path)
				if err != nil {
					t.Errorf("Error deleting temp file %s: %v", tt.path, err)
				}
			}()

			// check that running with the scenario it returns a different user
			output := strings.Builder{}
			config := smc.SyscallMonkeyConfig{
				ConfigPath:   fmt.Sprintf("./examples/%s", tt.scenario),
				TrailingArgs: []string{"whoami"},
				Silent:       true,
				TraceeStdout: &output,
			}
			smrun.RunTracer(&config, nil)

			modifiedOutput := strings.TrimSpace(output.String())

			if modifiedOutput != expectedUser {
				t.Errorf("Expected '%s', got '%s' (%s)", expectedUser, modifiedOutput, output.String())
			}
		})
	}

}
