package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"hawkop/internal/config"
)

// orgCmd represents the org command
var orgCmd = &cobra.Command{
	Use:   "org",
	Short: "Manage organization settings",
	Long: `Manage organization-related settings and operations.
	
Use subcommands to list organizations, set default organization, or view current organization settings.`,
}

// orgSetCmd sets the default organization ID
var orgSetCmd = &cobra.Command{
	Use:   "set <org-id>",
	Short: "Set the default organization ID",
	Long: `Set the default organization ID that will be used for subsequent commands.
	
The organization ID will be stored in your configuration file and used as the default
for commands that require an organization context.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runOrgSet(args[0])
	},
}

// orgGetCmd gets the current default organization ID
var orgGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Show the current default organization ID",
	Long:  `Display the currently configured default organization ID.`,
	Run: func(cmd *cobra.Command, args []string) {
		runOrgGet()
	},
}

// orgClearCmd clears the default organization ID
var orgClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear the default organization ID",
	Long:  `Remove the default organization ID from your configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		runOrgClear()
	},
}

func init() {
	rootCmd.AddCommand(orgCmd)
	orgCmd.AddCommand(orgSetCmd)
	orgCmd.AddCommand(orgGetCmd)
	orgCmd.AddCommand(orgClearCmd)
}

func runOrgSet(orgID string) {
	// Load existing config
	cfg, err := config.Load()
	checkError(err)

	// Validate that we have credentials
	if !cfg.HasValidCredentials() {
		fmt.Println("❌ No API key configured. Please run 'hawkop init' first.")
		return
	}

	// Set organization ID
	cfg.SetOrgID(orgID)

	// Save configuration
	err = cfg.Save()
	checkError(err)

	fmt.Printf("✅ Default organization ID set to: %s\n", orgID)
}

func runOrgGet() {
	// Load existing config
	cfg, err := config.Load()
	checkError(err)

	if cfg.OrgID == "" {
		fmt.Println("No default organization ID configured.")
		fmt.Println("Use 'hawkop org set <org-id>' to set one.")
	} else {
		fmt.Printf("Default organization ID: %s\n", cfg.OrgID)
	}
}

func runOrgClear() {
	// Load existing config
	cfg, err := config.Load()
	checkError(err)

	if cfg.OrgID == "" {
		fmt.Println("No default organization ID is currently set.")
		return
	}

	// Clear organization ID
	cfg.SetOrgID("")

	// Save configuration
	err = cfg.Save()
	checkError(err)

	fmt.Println("✅ Default organization ID cleared.")
}