package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

	"hawkop/internal/version"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show hawkop version information",
	Long:  `Display version information for hawkop including build details.`,
	Run: func(cmd *cobra.Command, args []string) {
		format, _ := cmd.Flags().GetString("format")
		runVersion(format)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.Flags().StringP("format", "f", "text", "Output format (text|json)")
}

func runVersion(outputFormat string) {
	switch outputFormat {
	case "json":
		info := version.GetInfo()
		data, err := json.MarshalIndent(info, "", "  ")
		if err != nil {
			fmt.Printf("❌ Failed to format JSON: %v\n", err)
			return
		}
		fmt.Println(string(data))
	case "text":
		fmt.Println(version.GetDetailedVersion())
	default:
		fmt.Printf("❌ Unknown format: %s. Use 'text' or 'json'\n", outputFormat)
	}
}
