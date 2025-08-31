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

// scanCmd represents the scan command
var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Manage scan-related operations",
	Long: `Manage scan-related operations including listing scans and viewing scan details.
	
Use subcommands to list scans, view scan details, or analyze scan results.`,
}

// scanListCmd lists scans in an organization
var scanListCmd = &cobra.Command{
	Use:   "list",
	Short: "List scans in an organization",
	Long: `List all scans for applications in the specified organization.
	
By default, uses your configured default organization and shows scans sorted by 
timestamp in descending order (most recent first). You can filter by application
name/ID and environment.`,
	Run: func(cmd *cobra.Command, args []string) {
		format, _ := cmd.Flags().GetString("format")
		limit, _ := cmd.Flags().GetInt("limit")
		org, _ := cmd.Flags().GetString("org")
		app, _ := cmd.Flags().GetString("app")
		env, _ := cmd.Flags().GetString("env")
		status, _ := cmd.Flags().GetString("status")
		sortBy, _ := cmd.Flags().GetString("sort-by")
		sortDir, _ := cmd.Flags().GetString("sort-dir")
		pageSize, _ := cmd.Flags().GetInt("page-size")
		pageToken, _ := cmd.Flags().GetString("page-token")
		runScanList(format, limit, org, app, env, status, sortBy, sortDir, pageSize, pageToken)
	},
}

// scanGetCmd gets details for a specific scan
var scanGetCmd = &cobra.Command{
	Use:   "get <scan-id>",
	Short: "Get details for a specific scan",
	Long: `Get detailed information about a specific scan including metadata,
duration, URL count, and alert statistics.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		scanID := args[0]
		format, _ := cmd.Flags().GetString("format")
		view, _ := cmd.Flags().GetString("view")
		runScanGet(scanID, format, view)
	},
}

// scanAlertsCmd lists alerts for a specific scan
var scanAlertsCmd = &cobra.Command{
	Use:   "alerts <scan-id>",
	Short: "List alerts for a specific scan",
	Long: `List all security alerts/findings for a specific scan.
	
Shows vulnerability details including severity, plugin ID, description, and URI count.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		scanID := args[0]
		format, _ := cmd.Flags().GetString("format")
		severity, _ := cmd.Flags().GetString("severity")
		limit, _ := cmd.Flags().GetInt("limit")
		runScanAlerts(scanID, format, severity, limit)
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.AddCommand(scanListCmd)
	scanCmd.AddCommand(scanGetCmd)
	scanCmd.AddCommand(scanAlertsCmd)

	// Add flags for scan list command
	scanListCmd.Flags().StringP("format", "f", "table", "Output format (table|json)")
	scanListCmd.Flags().IntP("limit", "l", 0, "Limit number of results (0 = no limit)")
	scanListCmd.Flags().StringP("org", "o", "", "Organization ID (uses default if not specified)")
	scanListCmd.Flags().StringP("app", "a", "", "Filter by application name or ID")
	scanListCmd.Flags().StringP("env", "e", "", "Filter by environment")
	scanListCmd.Flags().StringP("status", "s", "", "Filter by scan status (STARTED|COMPLETED|ERROR)")
	scanListCmd.Flags().StringP("sort-by", "", "timestamp", "Sort by field (timestamp|application|env|status)")
	scanListCmd.Flags().StringP("sort-dir", "", "desc", "Sort direction (asc|desc)")
	scanListCmd.Flags().IntP("page-size", "", 0, "Page size for API requests (default 1000, max 1000)")
	scanListCmd.Flags().StringP("page-token", "", "", "Page token for pagination")

	// Add flags for scan get command
	scanGetCmd.Flags().StringP("format", "f", "table", "Output format (table|json)")
	scanGetCmd.Flags().StringP("view", "v", "overview", "View type (overview|stats)")

	// Add flags for scan alerts command
	scanAlertsCmd.Flags().StringP("format", "f", "table", "Output format (table|json)")
	scanAlertsCmd.Flags().StringP("severity", "s", "", "Filter by severity (High|Medium|Low|Info)")
	scanAlertsCmd.Flags().IntP("limit", "l", 0, "Limit number of results (0 = no limit)")
}

func runScanList(outputFormat string, limit int, orgID string, appFilter string, envFilter string, statusFilter string, sortBy string, sortDir string, pageSize int, pageToken string) {
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

	// Build pagination options - always use max page size to minimize API requests
	paginationOpts := &api.PaginationOptions{
		PageSize: 1000, // Always use maximum to minimize API calls
	}
	
	// Override page size if explicitly set (but still cap at max)
	if pageSize > 0 {
		if pageSize > 1000 {
			pageSize = 1000
		}
		paginationOpts.PageSize = pageSize
	}
	
	if pageToken != "" {
		paginationOpts.PageToken = pageToken
	}
	
	// Only add sorting if explicitly different from defaults and not empty
	if sortBy != "" && sortBy != "timestamp" {
		paginationOpts.SortField = sortBy
	}
	if sortDir != "" && sortDir != "desc" {
		paginationOpts.SortDir = sortDir
	}

	// Get organization scans
	scanResults, err := client.ListOrganizationScansWithOptions(orgID, paginationOpts)
	if err != nil {
		fmt.Printf("❌ Failed to list scans: %v\n", err)
		return
	}

	// Apply filters
	filteredResults := []api.ApplicationScanResult{}
	for _, result := range scanResults {
		// App filter
		if appFilter != "" {
			appFilterLower := strings.ToLower(appFilter)
			if !strings.Contains(strings.ToLower(result.Scan.ApplicationName), appFilterLower) &&
			   !strings.Contains(strings.ToLower(result.Scan.ApplicationID), appFilterLower) {
				continue
			}
		}

		// Environment filter
		if envFilter != "" && !strings.EqualFold(result.Scan.Env, envFilter) {
			continue
		}

		// Status filter
		if statusFilter != "" && !strings.EqualFold(result.Scan.Status, statusFilter) {
			continue
		}

		filteredResults = append(filteredResults, result)
	}

	// Apply limit if specified
	if limit > 0 && len(filteredResults) > limit {
		filteredResults = filteredResults[:limit]
	}

	// Output based on format
	switch strings.ToLower(outputFormat) {
	case "json":
		outputScansJSON(filteredResults)
	case "table":
		outputScansTable(filteredResults)
	default:
		fmt.Printf("❌ Unknown format: %s. Use 'table' or 'json'\n", outputFormat)
		return
	}
}

func runScanGet(scanID string, outputFormat string, view string) {
	// This will need the specific scan details - for now we'll search through all scans
	cfg, err := config.Load()
	checkError(err)

	if !cfg.HasValidCredentials() {
		fmt.Println("❌ No API key configured. Please run 'hawkop init' first.")
		return
	}

	orgID := cfg.OrgID
	if orgID == "" {
		fmt.Println("❌ No organization configured. Set a default with 'hawkop org set <org-id>'")
		return
	}

	client := api.NewClient(cfg)
	scanResults, err := client.ListOrganizationScans(orgID)
	if err != nil {
		fmt.Printf("❌ Failed to get scan: %v\n", err)
		return
	}

	// Find the specific scan
	var targetScan *api.ApplicationScanResult
	for _, result := range scanResults {
		if result.Scan.ID == scanID {
			targetScan = &result
			break
		}
	}

	if targetScan == nil {
		fmt.Printf("❌ Scan not found: %s\n", scanID)
		return
	}

	// Output based on format and view
	switch strings.ToLower(outputFormat) {
	case "json":
		data, err := json.MarshalIndent(targetScan, "", "  ")
		if err != nil {
			fmt.Printf("❌ Failed to format JSON: %v\n", err)
			return
		}
		fmt.Println(string(data))
	case "table":
		outputScanDetailsTable(*targetScan, view)
	default:
		fmt.Printf("❌ Unknown format: %s. Use 'table' or 'json'\n", outputFormat)
	}
}

func runScanAlerts(scanID string, outputFormat string, severityFilter string, limit int) {
	cfg, err := config.Load()
	checkError(err)

	if !cfg.HasValidCredentials() {
		fmt.Println("❌ No API key configured. Please run 'hawkop init' first.")
		return
	}

	client := api.NewClient(cfg)
	alerts, err := client.GetScanAlerts(scanID)
	if err != nil {
		fmt.Printf("❌ Failed to get scan alerts: %v\n", err)
		return
	}

	// Apply severity filter if specified
	if severityFilter != "" {
		filteredAlerts := []api.ScanAlert{}
		for _, alert := range alerts {
			if strings.EqualFold(alert.Severity, severityFilter) {
				filteredAlerts = append(filteredAlerts, alert)
			}
		}
		alerts = filteredAlerts
	}

	// Apply limit if specified
	if limit > 0 && len(alerts) > limit {
		alerts = alerts[:limit]
	}

	// Output based on format
	switch strings.ToLower(outputFormat) {
	case "json":
		outputAlertsJSON(alerts)
	case "table":
		outputAlertsTable(alerts)
	default:
		fmt.Printf("❌ Unknown format: %s. Use 'table' or 'json'\n", outputFormat)
	}
}

func outputScansJSON(scanResults []api.ApplicationScanResult) {
	data, err := json.MarshalIndent(scanResults, "", "  ")
	if err != nil {
		fmt.Printf("❌ Failed to format JSON: %v\n", err)
		return
	}
	fmt.Println(string(data))
}

func outputScansTable(scanResults []api.ApplicationScanResult) {
	if len(scanResults) == 0 {
		fmt.Println("No scans found.")
		return
	}

	table := format.NewTable("SCAN ID", "APPLICATION", "ENV", "STATUS", "DURATION", "ALERTS", "TIMESTAMP")
	
	for _, result := range scanResults {
		// Format duration
		duration := ""
		if result.ScanDuration != nil {
			switch v := result.ScanDuration.(type) {
			case float64:
				duration = fmt.Sprintf("%.0fs", v)
			case string:
				if d, err := strconv.ParseFloat(v, 64); err == nil {
					duration = fmt.Sprintf("%.0fs", d)
				} else {
					duration = v
				}
			}
		}


		// Format alert count
		alertCount := ""
		if result.AlertStats != nil {
			alertCount = fmt.Sprintf("%d", result.AlertStats.Total)
		}

		// Format timestamp
		timestamp := ""
		if result.Scan.Timestamp != "" {
			if ts, err := strconv.ParseInt(result.Scan.Timestamp, 10, 64); err == nil {
				timestamp = time.Unix(ts/1000, 0).Format("2006-01-02 15:04")
			}
		}

		// Clean up values
		appName := result.Scan.ApplicationName
		if appName == "" {
			appName = "N/A"
		}

		env := result.Scan.Env
		if env == "" {
			env = "N/A"
		}

		status := result.Scan.Status
		if status == "" {
			status = "N/A"
		}

		table.AddRow(result.Scan.ID, appName, env, status, duration, alertCount, timestamp)
	}

	fmt.Print(table.Render())
}

func outputScanDetailsTable(scanResult api.ApplicationScanResult, view string) {
	switch view {
	case "overview":
		table := format.NewTable("FIELD", "VALUE")
		table.AddRow("Scan ID", scanResult.Scan.ID)
		table.AddRow("Application", scanResult.Scan.ApplicationName)
		table.AddRow("Environment", scanResult.Scan.Env)
		table.AddRow("Status", scanResult.Scan.Status)
		
		if scanResult.ScanDuration != nil {
			switch v := scanResult.ScanDuration.(type) {
			case float64:
				table.AddRow("Duration", fmt.Sprintf("%.0fs", v))
			case string:
				if d, err := strconv.ParseFloat(v, 64); err == nil {
					table.AddRow("Duration", fmt.Sprintf("%.0fs", d))
				} else {
					table.AddRow("Duration", v)
				}
			}
		}
		if scanResult.URLCount != nil {
			switch v := scanResult.URLCount.(type) {
			case float64:
				table.AddRow("URLs Scanned", fmt.Sprintf("%.0f", v))
			case string:
				table.AddRow("URLs Scanned", v)
			}
		}
		if scanResult.PolicyName != "" {
			table.AddRow("Policy", scanResult.PolicyName)
		}
		
		// Format timestamp
		if scanResult.Scan.Timestamp != "" {
			if ts, err := strconv.ParseInt(scanResult.Scan.Timestamp, 10, 64); err == nil {
				timestamp := time.Unix(ts/1000, 0).Format("2006-01-02 15:04:05")
				table.AddRow("Timestamp", timestamp)
			}
		}

		fmt.Print(table.Render())

	case "stats":
		if scanResult.AlertStats != nil {
			table := format.NewTable("SEVERITY", "COUNT")
			table.AddRow("High", fmt.Sprintf("%d", scanResult.AlertStats.High))
			table.AddRow("Medium", fmt.Sprintf("%d", scanResult.AlertStats.Medium))
			table.AddRow("Low", fmt.Sprintf("%d", scanResult.AlertStats.Low))
			table.AddRow("Info", fmt.Sprintf("%d", scanResult.AlertStats.Info))
			table.AddRow("Total", fmt.Sprintf("%d", scanResult.AlertStats.Total))
			fmt.Print(table.Render())
		} else {
			fmt.Println("No alert statistics available for this scan.")
		}

	default:
		fmt.Printf("❌ Unknown view: %s. Use 'overview' or 'stats'\n", view)
	}
}

func outputAlertsJSON(alerts []api.ScanAlert) {
	data, err := json.MarshalIndent(alerts, "", "  ")
	if err != nil {
		fmt.Printf("❌ Failed to format JSON: %v\n", err)
		return
	}
	fmt.Println(string(data))
}

func outputAlertsTable(alerts []api.ScanAlert) {
	if len(alerts) == 0 {
		fmt.Println("No alerts found.")
		return
	}

	table := format.NewTable("PLUGIN ID", "NAME", "SEVERITY", "URIS", "CWE")
	
	for _, alert := range alerts {
		// Clean up values
		name := alert.Name
		if name == "" {
			name = "N/A"
		}

		severity := alert.Severity
		if severity == "" {
			severity = "N/A"
		}

		uriCount := ""
		if alert.URICount > 0 {
			uriCount = fmt.Sprintf("%d", alert.URICount)
		} else {
			uriCount = "0"
		}

		cwe := alert.CWEID
		if cwe == "" {
			cwe = "N/A"
		}

		table.AddRow(alert.PluginID, name, severity, uriCount, cwe)
	}

	fmt.Print(table.Render())
}