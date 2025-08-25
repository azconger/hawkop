package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show hawkop version information",
	Long:  `Display version information for hawkop including build details.`,
	Run: func(cmd *cobra.Command, args []string) {
		runVersion()
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func runVersion() {
	fmt.Printf("hawkop version %s\n", Version)
	fmt.Printf("Built: %s\n", Date)
	fmt.Printf("Commit: %s\n", Commit)
	fmt.Printf("Go: %s\n", runtime.Version())
	fmt.Printf("Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
}