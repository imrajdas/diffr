package root

import (
	"github.com/imrajdas/diffr/pkg/cmd/version"
	"github.com/imrajdas/diffr/pkg/diffr"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "diffr [dir1/file1] [dir2/file2]",
	Example: "diffr /path/to/dir1 /path/to/dir2",
	Short:   "A web-based content difference analyzer",
	Long:    `A web-based tool to compare content differences between two directories/files` + "\n" + `Find more information at: https://github.com/imrajdas/diffr`,
	Args:    cobra.ExactArgs(2),
	Run:     diffr.RunWebServer,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().IntVarP(&diffr.Port, "port", "p", 8675, "Set the port for the web server to listen on, default is 8080")
	rootCmd.Flags().StringVarP(&diffr.Address, "address", "a", "http://localhost", "Set the address for the web server to listen on, default is http://localhost")

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(version.VersionCmd)
}
