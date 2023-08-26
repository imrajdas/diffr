package version

import (
	"os"

	"github.com/spf13/cobra"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Displays the version of diffr",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		version := os.Getenv("VERSION")
		if version == "" {
			version = "develop"
		}

		cmd.Printf("Diffr Version: %s", version)
	},
}
