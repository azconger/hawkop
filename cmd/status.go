package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"hawkop/internal/config"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show hawkop configuration and connection status",
	Long: `Display the current hawkop configuration including:
- Connection status
- API key status
- Default organization
- JWT token status
- Configuration file location`,
	Run: func(cmd *cobra.Command, args []string) {
		runStatus()
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}

func runStatus() {
	fmt.Println("🦅 HawkOp Status")
	fmt.Println("================")
	fmt.Println()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("❌ Configuration Error: %v\n", err)
		return
	}

	// Display configuration file location
	fmt.Printf("📁 Config file: %s\n", config.GetConfigFile())
	fmt.Println()

	// Check API key status
	if cfg.APIKey == "" {
		fmt.Println("🔑 API Key: ❌ Not configured")
		fmt.Println("   Run 'hawkop init' to set up your API key")
	} else {
		fmt.Println("🔑 API Key: ✅ Configured")
		fmt.Printf("   Key: %s...%s\n", 
			cfg.APIKey[:min(8, len(cfg.APIKey))], 
			strings.Repeat("*", max(0, len(cfg.APIKey)-8)))
	}
	fmt.Println()

	// Check organization status
	if cfg.OrgID == "" {
		fmt.Println("🏢 Default Org: ❌ Not set")
		fmt.Println("   Use 'hawkop org set <org-id>' to set a default organization")
	} else {
		fmt.Println("🏢 Default Org: ✅ Set")
		fmt.Printf("   Organization ID: %s\n", cfg.OrgID)
	}
	fmt.Println()

	// Check JWT status
	if cfg.JWT == nil {
		fmt.Println("🎫 JWT Token: ❌ None")
		if cfg.HasValidCredentials() {
			fmt.Println("   A token will be automatically obtained when needed")
		}
	} else if cfg.JWT.IsExpired() {
		fmt.Println("🎫 JWT Token: ⏰ Expired")
		fmt.Printf("   Expired at: %s\n", cfg.JWT.ExpiresAt.Format("2006-01-02 15:04:05 MST"))
		fmt.Println("   A fresh token will be obtained automatically")
	} else {
		fmt.Println("🎫 JWT Token: ✅ Valid")
		fmt.Printf("   Expires at: %s\n", cfg.JWT.ExpiresAt.Format("2006-01-02 15:04:05 MST"))
	}
	fmt.Println()

	// Overall status
	if !cfg.HasValidCredentials() {
		fmt.Println("🔗 Overall Status: ❌ Not ready")
		fmt.Println("   Please run 'hawkop init' to configure your API key")
	} else {
		fmt.Println("🔗 Overall Status: ✅ Ready")
		fmt.Println("   You can now use hawkop commands")
	}
}