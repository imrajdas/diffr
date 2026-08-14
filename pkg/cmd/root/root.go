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
	Run:     diffr.Run,
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().IntVarP(&diffr.Port, "port", "p", 8675, "Set the port for the web server to listen on")
	rootCmd.Flags().StringVarP(&diffr.Address, "address", "a", "http://localhost", "Bind address for the web server (host or URL)")
	rootCmd.Flags().IntVarP(&diffr.Context, "context", "C", 3, "Number of unified-diff context lines")
	rootCmd.Flags().BoolVarP(&diffr.IgnoreWhitespace, "ignore-whitespace", "w", false, "Do not report files that differ only in whitespace")
	rootCmd.Flags().BoolVarP(&diffr.IgnoreCase, "ignore-case", "i", false, "Do not report files that differ only in case")
	rootCmd.Flags().BoolVar(&diffr.Stdout, "stdout", false, "Print diffs to stdout instead of starting the web UI")
	rootCmd.Flags().BoolVar(&diffr.JSON, "json", false, "Print a JSON report to stdout instead of starting the web UI")
	rootCmd.Flags().BoolVar(&diffr.NoOpen, "no-open", false, "Do not open a browser when starting the web UI")
	rootCmd.Flags().StringVar(&diffr.PatchFile, "patch", "", "Write a unified patch to the given file instead of starting the web UI (use - for stdout)")
	rootCmd.Flags().StringArrayVarP(&diffr.Exclude, "exclude", "e", nil, "Gitignore-style glob to exclude (repeatable)")
	rootCmd.Flags().StringVar(&diffr.IgnoreFile, "ignore-file", "", "Path to a .diffrignore file (gitignore syntax)")
	rootCmd.Flags().BoolVar(&diffr.NoDefaultExclude, "no-default-exclude", false, "Do not apply built-in excludes (.git, node_modules, ...)")
	rootCmd.Flags().BoolVar(&diffr.NoGitignore, "no-gitignore", false, "Do not apply .gitignore from git repositories")

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.AddCommand(version.VersionCmd)
}
