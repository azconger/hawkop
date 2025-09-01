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

// teamCmd represents the team command
var teamCmd = &cobra.Command{
	Use:   "team",
	Short: "Manage team-related operations",
	Long: `Manage team-related operations including listing teams in organizations.
	
Use subcommands to list teams, view team details, or manage team settings.`,
}

// teamListCmd lists teams in an organization
var teamListCmd = &cobra.Command{
	Use:   "list",
	Short: "List teams in an organization",
	Long: `List all teams that belong to the specified organization.
	
By default, uses your configured default organization. You can specify a different
organization using the --org flag. This command requires ADMIN or OWNER role.`,
	Run: func(cmd *cobra.Command, args []string) {
		format, _ := cmd.Flags().GetString("format")
		limit, _ := cmd.Flags().GetInt("limit")
		org, _ := cmd.Flags().GetString("org")
		runTeamList(format, limit, org)
	},
}

func init() {
	rootCmd.AddCommand(teamCmd)
	teamCmd.AddCommand(teamListCmd)

	// Add flags for team list command
	teamListCmd.Flags().StringP("format", "f", "table", "Output format (table|json)")
	teamListCmd.Flags().IntP("limit", "l", 0, "Limit number of results (0 = no limit)")
	teamListCmd.Flags().StringP("org", "o", "", "Organization ID (uses default if not specified)")
}

func runTeamList(outputFormat string, limit int, orgID string) {
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

	// Get organization teams
	teams, err := client.ListOrganizationTeams(orgID)
	if err != nil {
		fmt.Printf("❌ Failed to list teams: %v\n", err)
		return
	}

	// Apply limit if specified
	if limit > 0 && len(teams) > limit {
		teams = teams[:limit]
	}

	// Output based on format
	switch strings.ToLower(outputFormat) {
	case "json":
		outputTeamsJSON(teams)
	case "table":
		outputTeamsTable(teams)
	default:
		fmt.Printf("❌ Unknown format: %s. Use 'table' or 'json'\n", outputFormat)
		return
	}
}

func outputTeamsJSON(teams []api.Team) {
	data, err := json.MarshalIndent(teams, "", "  ")
	if err != nil {
		fmt.Printf("❌ Failed to format JSON: %v\n", err)
		return
	}
	fmt.Println(string(data))
}

func outputTeamsTable(teams []api.Team) {
	if len(teams) == 0 {
		fmt.Println("No teams found.")
		return
	}

	table := format.NewTable("ID", "NAME", "USERS", "APPS", "CREATED")

	for _, team := range teams {
		// Count users and applications
		userCount := fmt.Sprintf("%d", len(team.Users))
		appCount := fmt.Sprintf("%d", len(team.Applications))

		// Format created date
		created := ""
		if team.CreatedTimestamp != "" {
			if ts, err := strconv.ParseInt(team.CreatedTimestamp, 10, 64); err == nil {
				created = time.Unix(ts/1000, 0).Format("2006-01-02")
			}
		}

		// Clean up values
		name := team.Name
		if name == "" {
			name = "N/A"
		}

		table.AddRow(team.ID, name, userCount, appCount, created)
	}

	fmt.Print(table.Render())
}
