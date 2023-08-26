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
		cmd.Printf("Diffr Version: %s", os.Getenv("VERSION"))
	},
}
