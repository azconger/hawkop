package cmd

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"hawkop/internal/api"
	"hawkop/internal/config"
)

type ScanCommandTestSuite struct {
	suite.Suite
	mockClient *api.MockClient
	origOutput io.Writer
	buffer     *bytes.Buffer
}

func (suite *ScanCommandTestSuite) SetupTest() {
	suite.mockClient = api.NewMockClient()

	// Capture output
	suite.buffer = new(bytes.Buffer)
	suite.origOutput = os.Stdout
}

func (suite *ScanCommandTestSuite) TearDownTest() {
	// Reset output if needed
}

func (suite *ScanCommandTestSuite) TestScanListCommand_Success() {
	// Mock successful API response
	mockScans := []api.ApplicationScanResult{
		{
			Scan: api.Scan{
				ID:              "scan-1",
				ApplicationID:   "app-1",
				ApplicationName: "Test App",
				Status:          "COMPLETED",
				Timestamp:       "1756596062834",
				Env:             "production",
			},
			ScanDuration: "45",
			URLCount:     "10",
			AlertStats: &api.AlertStats{
				High:   2,
				Medium: 3,
				Low:    1,
				Total:  6,
			},
		},
	}

	suite.mockClient.On("ListOrganizationScans", "test-org-id").Return(mockScans, nil)

	// Test the command logic (this would require refactoring to inject the client)
	// For now, we'll test the command structure
	cmd := scanListCmd
	assert.Equal(suite.T(), "list", cmd.Use)
	assert.Contains(suite.T(), cmd.Short, "List scans")
}

func (suite *ScanCommandTestSuite) TestScanGetCommand_Success() {
	// Mock successful API response
	mockScans := []api.ApplicationScanResult{
		{
			Scan: api.Scan{
				ID:              "scan-1",
				ApplicationID:   "app-1",
				ApplicationName: "Test App",
				Status:          "COMPLETED",
				Timestamp:       "1756596062834",
			},
		},
	}

	suite.mockClient.On("ListOrganizationScans", "test-org-id").Return(mockScans, nil)

	cmd := scanGetCmd
	assert.Equal(suite.T(), "get <scan-id>", cmd.Use)
	assert.Contains(suite.T(), cmd.Short, "Get details")
}

func (suite *ScanCommandTestSuite) TestScanAlertsCommand_Success() {
	// Mock successful alerts response
	mockAlerts := []api.ScanAlert{
		{
			PluginID: "10001",
			Name:     "SQL Injection",
			Severity: "High",
			URICount: 3,
			CWEID:    "CWE-89",
		},
	}

	suite.mockClient.On("GetScanAlerts", "scan-1").Return(mockAlerts, nil)

	cmd := scanAlertsCmd
	assert.Equal(suite.T(), "alerts <scan-id>", cmd.Use)
	assert.Contains(suite.T(), cmd.Short, "List alerts")
}

func (suite *ScanCommandTestSuite) TestScanCommand_Structure() {
	// Test main scan command structure
	assert.Equal(suite.T(), "scan", scanCmd.Use)
	assert.Contains(suite.T(), scanCmd.Short, "Manage scan-related operations")

	// Verify subcommands are registered
	subcommands := []string{}
	for _, cmd := range scanCmd.Commands() {
		subcommands = append(subcommands, cmd.Use)
	}

	assert.Contains(suite.T(), subcommands, "list")
	assert.Contains(suite.T(), subcommands, "get <scan-id>")
	assert.Contains(suite.T(), subcommands, "alerts <scan-id>")
}

func (suite *ScanCommandTestSuite) TestScanListFlags() {
	cmd := scanListCmd

	// Test that all expected flags are present
	formatFlag := cmd.Flags().Lookup("format")
	assert.NotNil(suite.T(), formatFlag)
	assert.Equal(suite.T(), "table", formatFlag.DefValue)

	limitFlag := cmd.Flags().Lookup("limit")
	assert.NotNil(suite.T(), limitFlag)
	assert.Equal(suite.T(), "0", limitFlag.DefValue)

	orgFlag := cmd.Flags().Lookup("org")
	assert.NotNil(suite.T(), orgFlag)

	appFlag := cmd.Flags().Lookup("app")
	assert.NotNil(suite.T(), appFlag)

	envFlag := cmd.Flags().Lookup("env")
	assert.NotNil(suite.T(), envFlag)

	statusFlag := cmd.Flags().Lookup("status")
	assert.NotNil(suite.T(), statusFlag)
}

func (suite *ScanCommandTestSuite) TestScanGetFlags() {
	cmd := scanGetCmd

	formatFlag := cmd.Flags().Lookup("format")
	assert.NotNil(suite.T(), formatFlag)
	assert.Equal(suite.T(), "table", formatFlag.DefValue)

	viewFlag := cmd.Flags().Lookup("view")
	assert.NotNil(suite.T(), viewFlag)
	assert.Equal(suite.T(), "overview", viewFlag.DefValue)
}

func (suite *ScanCommandTestSuite) TestScanAlertsFlags() {
	cmd := scanAlertsCmd

	formatFlag := cmd.Flags().Lookup("format")
	assert.NotNil(suite.T(), formatFlag)
	assert.Equal(suite.T(), "table", formatFlag.DefValue)

	severityFlag := cmd.Flags().Lookup("severity")
	assert.NotNil(suite.T(), severityFlag)

	limitFlag := cmd.Flags().Lookup("limit")
	assert.NotNil(suite.T(), limitFlag)
	assert.Equal(suite.T(), "0", limitFlag.DefValue)
}

// Helper function to test command execution with mock config
func executeCommandWithMockConfig(cmd *cobra.Command, args []string, cfg *config.Config) error {
	// This would require dependency injection to properly test
	// For now, we test the command structure and flags
	cmd.SetArgs(args)
	return cmd.Execute()
}

func TestScanCommandTestSuite(t *testing.T) {
	suite.Run(t, new(ScanCommandTestSuite))
}
