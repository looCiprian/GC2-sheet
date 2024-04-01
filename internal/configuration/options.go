package configuration

import "net/url"

type options struct {
	credential string
	sheetId    string
	driveId    string
	proxy      *url.URL
	debug      bool
}

type ConfigurationFile struct {
	Key     string `yaml:"key"`
	Sheet   string `yaml:"sheet"`
	Drive   string `yaml:"drive"`
	Proxy   string `yaml:"proxy"`
	Verbose bool   `yaml:"verbose" default:"false"`
}

var command options

func SetOptions(credential string, sheetId string, driveId string, proxy *url.URL, debug bool) {

	command.credential = credential
	command.sheetId = sheetId
	command.driveId = driveId
	command.proxy = proxy
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

func GetOptionsProxy() *url.URL {

	return command.proxy

}

func GetOptionsDebug() bool {

	return command.debug

}
