package cmd

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"hawkop/internal/api"
	"hawkop/internal/config"
	"hawkop/internal/format"
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

// orgListCmd lists all organizations the user belongs to
var orgListCmd = &cobra.Command{
	Use:   "list",
	Short: "List organizations you belong to",
	Long: `List all organizations that you have access to in StackHawk.
	
This command displays your organization memberships including organization ID, 
name, plan, and other details.`,
	Run: func(cmd *cobra.Command, args []string) {
		format, _ := cmd.Flags().GetString("format")
		limit, _ := cmd.Flags().GetInt("limit")
		runOrgList(format, limit)
	},
}

func init() {
	rootCmd.AddCommand(orgCmd)
	orgCmd.AddCommand(orgSetCmd)
	orgCmd.AddCommand(orgGetCmd)
	orgCmd.AddCommand(orgClearCmd)
	orgCmd.AddCommand(orgListCmd)

	// Add flags for org list command
	orgListCmd.Flags().StringP("format", "f", "table", "Output format (table|json)")
	orgListCmd.Flags().IntP("limit", "l", 0, "Limit number of results (0 = no limit)")
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

func runOrgList(outputFormat string, limit int) {
	// Load configuration
	cfg, err := config.Load()
	checkError(err)

	// Validate that we have credentials
	if !cfg.HasValidCredentials() {
		fmt.Println("❌ No API key configured. Please run 'hawkop init' first.")
		return
	}

	// Create API client
	client := api.NewClient(cfg)

	// Get organizations
	orgs, err := client.ListOrganizations()
	if err != nil {
		fmt.Printf("❌ Failed to list organizations: %v\n", err)
		return
	}

	// Apply limit if specified
	if limit > 0 && len(orgs) > limit {
		orgs = orgs[:limit]
	}

	// Output based on format
	switch strings.ToLower(outputFormat) {
	case "json":
		outputJSON(orgs)
	case "table":
		outputTable(orgs)
	default:
		fmt.Printf("❌ Unknown format: %s. Use 'table' or 'json'\n", outputFormat)
		return
	}
}

func outputJSON(orgs []api.Organization) {
	data, err := json.MarshalIndent(orgs, "", "  ")
	if err != nil {
		fmt.Printf("❌ Failed to format JSON: %v\n", err)
		return
	}
	fmt.Println(string(data))
}

func outputTable(orgs []api.Organization) {
	if len(orgs) == 0 {
		fmt.Println("No organizations found.")
		return
	}

	table := format.NewTable("ID", "NAME", "PLAN", "CREATED")

	for _, org := range orgs {
		created := ""
		if org.CreatedTimestamp != "" {
			// Convert millisecond timestamp to readable date
			if ts, err := strconv.ParseInt(org.CreatedTimestamp, 10, 64); err == nil {
				created = time.Unix(ts/1000, 0).Format("2006-01-02")
			}
		}

		plan := org.Plan
		if plan == "" {
			plan = "N/A"
		}

		table.AddRow(org.ID, org.Name, plan, created)
	}

	fmt.Print(table.Render())
}
