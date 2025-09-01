package cmd

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"hawkop/internal/api"
	"hawkop/internal/config"
	"hawkop/internal/format"
)

// appCmd represents the app command
var appCmd = &cobra.Command{
	Use:   "app",
	Short: "Manage application-related operations",
	Long: `Manage application-related operations including listing applications in organizations.
	
Use subcommands to list applications, view application details, or manage application settings.`,
}

// appListCmd lists applications in an organization
var appListCmd = &cobra.Command{
	Use:   "list",
	Short: "List applications in an organization",
	Long: `List all applications that belong to the specified organization.
	
By default, uses your configured default organization. You can specify a different
organization using the --org flag. This command requires appropriate permissions.`,
	Run: func(cmd *cobra.Command, args []string) {
		format, _ := cmd.Flags().GetString("format")
		limit, _ := cmd.Flags().GetInt("limit")
		org, _ := cmd.Flags().GetString("org")
		status, _ := cmd.Flags().GetString("status")
		runAppList(format, limit, org, status)
	},
}

func init() {
	rootCmd.AddCommand(appCmd)
	appCmd.AddCommand(appListCmd)

	// Add flags for app list command
	appListCmd.Flags().StringP("format", "f", "table", "Output format (table|json)")
	appListCmd.Flags().IntP("limit", "l", 0, "Limit number of results (0 = no limit)")
	appListCmd.Flags().StringP("org", "o", "", "Organization ID (uses default if not specified)")
	appListCmd.Flags().StringP("status", "s", "", "Filter by application status (ACTIVE|ENV_INCOMPLETE)")
}

func runAppList(outputFormat string, limit int, orgID string, statusFilter string) {
	// Load configuration
	cfg, err := config.Load()
	checkError(err)

	// Validate that we have credentials
	if !cfg.HasValidCredentials() {
		fmt.Println("❌ No API key configured. Please run 'hawkop init' first.")
		return
	}

	// Determine which organization to use
	if orgID == "" {
		orgID = cfg.OrgID
		if orgID == "" {
			fmt.Println("❌ No organization specified. Use --org flag or set a default with 'hawkop org set <org-id>'")
			return
		}
	}

	// Create API client
	client := api.NewClient(cfg)

	// Get organization applications
	applications, err := client.ListOrganizationApplications(orgID)
	if err != nil {
		fmt.Printf("❌ Failed to list applications: %v\n", err)
		return
	}

	// Apply status filter if specified
	if statusFilter != "" {
		filteredApps := []api.AppApplication{}
		statusFilterUpper := strings.ToUpper(statusFilter)
		for _, app := range applications {
			if strings.ToUpper(app.ApplicationStatus) == statusFilterUpper {
				filteredApps = append(filteredApps, app)
			}
		}
		applications = filteredApps
	}

	// Apply limit if specified
	if limit > 0 && len(applications) > limit {
		applications = applications[:limit]
	}

	// Output based on format
	switch strings.ToLower(outputFormat) {
	case "json":
		outputApplicationsJSON(applications)
	case "table":
		outputApplicationsTable(applications)
	default:
		fmt.Printf("❌ Unknown format: %s. Use 'table' or 'json'\n", outputFormat)
		return
	}
}

func outputApplicationsJSON(applications []api.AppApplication) {
	data, err := json.MarshalIndent(applications, "", "  ")
	if err != nil {
		fmt.Printf("❌ Failed to format JSON: %v\n", err)
		return
	}
	fmt.Println(string(data))
}

func outputApplicationsTable(applications []api.AppApplication) {
	if len(applications) == 0 {
		fmt.Println("No applications found.")
		return
	}

	table := format.NewTable("ID", "NAME", "ENV", "STATUS", "TYPE")

	for _, app := range applications {
		// Clean up values
		name := app.Name
		if name == "" {
			name = "N/A"
		}

		env := app.Env
		if env == "" {
			env = "N/A"
		}

		status := app.ApplicationStatus
		if status == "" {
			status = "N/A"
		}

		appType := app.ApplicationType
		if appType == "" {
			appType = "N/A"
		}

		table.AddRow(app.ApplicationID, name, env, status, appType)
	}

	fmt.Print(table.Render())
}
