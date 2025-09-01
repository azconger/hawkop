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

// userCmd represents the user command
var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage user-related operations",
	Long: `Manage user-related operations including listing users in organizations.
	
Use subcommands to list users, view user details, or manage user settings.`,
}

// userListCmd lists users in an organization
var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "List users in an organization",
	Long: `List all users that belong to the specified organization.
	
By default, uses your configured default organization. You can specify a different
organization using the --org flag. This command requires ADMIN or OWNER role.`,
	Run: func(cmd *cobra.Command, args []string) {
		format, _ := cmd.Flags().GetString("format")
		limit, _ := cmd.Flags().GetInt("limit")
		org, _ := cmd.Flags().GetString("org")
		role, _ := cmd.Flags().GetString("role")
		runUserList(format, limit, org, role)
	},
}

func init() {
	rootCmd.AddCommand(userCmd)
	userCmd.AddCommand(userListCmd)

	// Add flags for user list command
	userListCmd.Flags().StringP("format", "f", "table", "Output format (table|json)")
	userListCmd.Flags().IntP("limit", "l", 0, "Limit number of results (0 = no limit)")
	userListCmd.Flags().StringP("org", "o", "", "Organization ID (uses default if not specified)")
	userListCmd.Flags().StringP("role", "r", "", "Filter by user role (admin|member|owner)")
}

func runUserList(outputFormat string, limit int, orgID string, roleFilter string) {
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

	// Get organization members
	members, err := client.ListOrganizationMembers(orgID)
	if err != nil {
		fmt.Printf("❌ Failed to list users: %v\n", err)
		return
	}

	// Apply role filter if specified
	if roleFilter != "" {
		filteredMembers := []api.OrganizationMember{}
		roleFilterUpper := strings.ToUpper(roleFilter)
		for _, member := range members {
			// Extract role from organizations array to check against filter
			memberRole := ""
			if member.External != nil {
				for _, orgMembership := range member.External.Organizations {
					memberRole = orgMembership.Role
					break
				}
			}
			if strings.ToUpper(memberRole) == roleFilterUpper {
				filteredMembers = append(filteredMembers, member)
			}
		}
		members = filteredMembers
	}

	// Apply limit if specified
	if limit > 0 && len(members) > limit {
		members = members[:limit]
	}

	// Output based on format
	switch strings.ToLower(outputFormat) {
	case "json":
		outputUsersJSON(members)
	case "table":
		outputUsersTable(members)
	default:
		fmt.Printf("❌ Unknown format: %s. Use 'table' or 'json'\n", outputFormat)
		return
	}
}

func outputUsersJSON(members []api.OrganizationMember) {
	data, err := json.MarshalIndent(members, "", "  ")
	if err != nil {
		fmt.Printf("❌ Failed to format JSON: %v\n", err)
		return
	}
	fmt.Println(string(data))
}

func outputUsersTable(members []api.OrganizationMember) {
	if len(members) == 0 {
		fmt.Println("No users found.")
		return
	}

	table := format.NewTable("NAME", "EMAIL", "ROLE", "PROVIDER", "CREATED")

	for _, member := range members {
		name := ""
		email := ""
		role := ""

		// Extract user info from External field
		if member.External != nil {
			name = member.External.FullName
			if name == "" {
				name = fmt.Sprintf("%s %s", member.External.FirstName, member.External.LastName)
			}
			email = member.External.Email

			// Extract role from the organizations array in External
			for _, orgMembership := range member.External.Organizations {
				role = orgMembership.Role
				break // Use the first organization role (should match the requested org)
			}
		}

		// Format provider
		provider := ""
		if member.Provider != nil {
			provider = member.Provider.Slug
		}

		// Format created date
		created := ""
		if member.CreatedTimestamp != "" {
			if ts, err := strconv.ParseInt(member.CreatedTimestamp, 10, 64); err == nil {
				created = time.Unix(ts/1000, 0).Format("2006-01-02")
			}
		}

		// Clean up values
		if name == "" {
			name = "N/A"
		}
		if email == "" {
			email = "N/A"
		}
		if role == "" {
			role = "N/A"
		}
		if provider == "" {
			provider = "N/A"
		}

		table.AddRow(name, email, role, provider, created)
	}

	fmt.Print(table.Render())
}
