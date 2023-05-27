package configuration

type options struct {
	credential string
	sheetId    string
	driveId    string
	debug      bool
}

type ConfigurationFile struct {
	Key     string `yaml:"key"`
	Sheet   string `yaml:"sheet"`
	Drive   string `yaml:"drive"`
	Verbose bool   `yaml:"verbose" default:"false"`
}

var command options

func SetOptions(credential string, sheetId string, driveId string, debug bool) {

	command.credential = credential
	command.sheetId = sheetId
	command.driveId = driveId
	command.debug = debug

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

func GetOptionsDebug() bool {

	return command.debug

}
