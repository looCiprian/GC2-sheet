package configuration

import (
	"net/url"
)

type options struct {
	commandService          ConnectorService
	fileSystemService       ConnectorService
	googleServiceAccountKey string
	googleSheetID           string
	googleDriveID           string
	microsoftTenantID       string
	microsoftClientID       string
	microsoftClientSecret   string
	microsoftSiteID         string
	rowID                   int
	proxy                   *url.URL
	verbose                 bool
}

var command options

func SetOptions(
	commandService,
	fileSystemService,
	googleServiceAccountKey,
	googleSheetID,
	googleDriveID,
	microsoftTenantID,
	microsoftClientID,
	microsoftClientSecret,
	microsoftSiteID string,
	rowID int,
	proxy *url.URL,
	verbose bool,
) {
	switch commandService {
	case Google.String():
		command.commandService = Google
	case Microsoft.String():
		command.commandService = Microsoft
	}

	switch fileSystemService {
	case Google.String():
		command.fileSystemService = Google
	case Microsoft.String():
		command.fileSystemService = Microsoft
	}

	command.googleServiceAccountKey = googleServiceAccountKey
	command.googleSheetID = googleSheetID
	command.googleDriveID = googleDriveID
	command.microsoftTenantID = microsoftTenantID
	command.microsoftClientID = microsoftClientID
	command.microsoftClientSecret = microsoftClientSecret
	command.microsoftSiteID = microsoftSiteID
	command.proxy = proxy
	command.rowID = rowID
	command.verbose = verbose

}

func GetOptionsCommandService() ConnectorService { return command.commandService }

func GetOptionsFileSystemService() ConnectorService { return command.fileSystemService }

func NeedsGoogleConnectorService() bool {
	return command.commandService == Google || command.fileSystemService == Google
}

func NeedsMicrosoftConnectorService() bool {
	return command.commandService == Microsoft || command.fileSystemService == Microsoft
}

func GetOptionsGoogleServiceAccountKey() string {

	return command.googleServiceAccountKey

}

func GetOptionsGoogleSheetID() string {

	return command.googleSheetID

}

func GetOptionsGoogleDriveID() string {

	return command.googleDriveID

}

func GetOptionsProxy() *url.URL {

	return command.proxy

}

func GetOptionsDebug() bool {

	return command.verbose

}

func GetSourceFirstCommandIndex() int {
	return command.rowID
}

func GetOptionsMicrosoftTenantID() string {
	return command.microsoftTenantID
}
func GetOptionsMicrosoftClientID() string {
	return command.microsoftClientID
}
func GetOptionsMicrosoftClientSecret() string {
	return command.microsoftClientSecret
}
func GetOptionsMicrosoftSiteID() string {
	return command.microsoftSiteID
}
