package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
	"golang.org/x/term"
	"hawkop/internal/config"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize hawkop with your StackHawk API key",
	Long: `Initialize hawkop by setting up your StackHawk API key and optional default organization.
	
The API key will be securely stored in your local configuration file and used for
authenticating with the StackHawk API. You can optionally set a default organization
to use for subsequent commands.`,
	Run: func(cmd *cobra.Command, args []string) {
		runInit()
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit() {
	fmt.Println("🦅 Welcome to HawkOp!")
	fmt.Println()
	fmt.Println("Let's set up your StackHawk credentials...")
	fmt.Println()

	// Load existing config
	cfg, err := config.Load()
	checkError(err)

	// Prompt for API key
	apiKey, err := promptForAPIKey(cfg.APIKey)
	checkError(err)

	if apiKey != "" {
		cfg.SetAPIKey(apiKey)
	}

	// Prompt for default organization (optional)
	orgID, err := promptForOrgID(cfg.OrgID)
	checkError(err)

	if orgID != "" {
		cfg.SetOrgID(orgID)
	}

	// Save configuration
	err = cfg.Save()
	checkError(err)

	fmt.Println()
	fmt.Println("✅ Configuration saved successfully!")
	fmt.Printf("   Config file: %s\n", config.GetConfigFile())

	if cfg.APIKey != "" {
		fmt.Println("   API key: configured")
	}
	if cfg.OrgID != "" {
		fmt.Printf("   Default org ID: %s\n", cfg.OrgID)
	}

	fmt.Println()
	fmt.Println("You can now use hawkop commands. Try:")
	fmt.Println("  hawkop status")
	fmt.Println("  hawkop org list")
}

func promptForAPIKey(currentKey string) (string, error) {
	if currentKey != "" {
		fmt.Printf("Current API key: %s...%s\n",
			currentKey[:min(8, len(currentKey))],
			strings.Repeat("*", max(0, len(currentKey)-8)))
		fmt.Print("Enter new API key (or press Enter to keep current): ")
	} else {
		fmt.Print("Enter your StackHawk API key: ")
	}

	// Read password without echo
	byteKey, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", fmt.Errorf("failed to read API key: %w", err)
	}

	fmt.Println() // Print newline after hidden input

	apiKey := strings.TrimSpace(string(byteKey))

	// If empty and we have a current key, keep the current key
	if apiKey == "" && currentKey != "" {
		return "", nil // Return empty to indicate no change
	}

	if apiKey == "" {
		return "", fmt.Errorf("API key is required")
	}

	return apiKey, nil
}

func promptForOrgID(currentOrgID string) (string, error) {
	reader := bufio.NewReader(os.Stdin)

	if currentOrgID != "" {
		fmt.Printf("Current default org ID: %s\n", currentOrgID)
		fmt.Print("Enter new org ID (or press Enter to keep current): ")
	} else {
		fmt.Print("Enter default org ID (optional): ")
	}

	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read org ID: %w", err)
	}

	orgID := strings.TrimSpace(input)

	// If empty and we have a current org ID, keep the current org ID
	if orgID == "" && currentOrgID != "" {
		return "", nil // Return empty to indicate no change
	}

	return orgID, nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
