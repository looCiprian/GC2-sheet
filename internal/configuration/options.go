package configuration

type options struct {
	credential string
	sheetId    string
	driveId    string
}

var command options

func SetOptions(credential string, sheetId string, driveId string) {

	command.credential = credential
	command.sheetId = sheetId
	command.driveId = driveId

}

func GetOptionsCredential() string {

	return command.credential

}

func GetOptionsSheetId() string {

	return command.sheetId

}

func GetOptionsDriveId() string {

	return command.driveId

}