//go:build !windows

package C2

import (
	"os/exec"
	"strings"
)

func executeCommand(commandToExecute string) string {

	var arguments []string
	var outCommand []byte
	var err error

	splitArgs := strings.Split(commandToExecute, " ")
	if len(splitArgs) > 1 {
		// Get arguments after command Example: -la
		arguments = splitArgs[1:]
		// Get command without arguments Example: ls
		commandToExecute = splitArgs[0]
	}
	outCommand, err = exec.Command(commandToExecute, arguments...).Output()

	if err != nil {
		out := err.Error()
		return out
	}

	out := string(outCommand)
	return out
}
