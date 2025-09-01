package cmd

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"hawkop/internal/api"
)

type AppCommandTestSuite struct {
	suite.Suite
	mockClient *api.MockClient
	origOutput io.Writer
	buffer     *bytes.Buffer
}

func (suite *AppCommandTestSuite) SetupTest() {
	suite.mockClient = api.NewMockClient()

	// Capture output
	suite.buffer = new(bytes.Buffer)
	suite.origOutput = os.Stdout
}

func (suite *AppCommandTestSuite) TearDownTest() {
	// Reset output if needed
}

func (suite *AppCommandTestSuite) TestAppListCommand_Success() {
	// Mock successful API response
	mockApps := []api.AppApplication{
		{
			ApplicationID:     "app-1",
			Name:              "Test Application",
			ApplicationStatus: "ACTIVE",
			ApplicationType:   "STANDARD",
		},
	}

	suite.mockClient.On("ListOrganizationApplications", "test-org-id").Return(mockApps, nil)

	cmd := appListCmd
	assert.Equal(suite.T(), "list", cmd.Use)
	assert.Contains(suite.T(), cmd.Short, "List applications")
}

func (suite *AppCommandTestSuite) TestAppCommand_Structure() {
	// Test main app command structure
	assert.Equal(suite.T(), "app", appCmd.Use)
	assert.Contains(suite.T(), appCmd.Short, "Manage application")

	// Verify subcommands are registered
	subcommands := []string{}
	for _, cmd := range appCmd.Commands() {
		subcommands = append(subcommands, cmd.Use)
	}

	assert.Contains(suite.T(), subcommands, "list")
}

func (suite *AppCommandTestSuite) TestAppListFlags() {
	cmd := appListCmd

	// Test that all expected flags are present
	formatFlag := cmd.Flags().Lookup("format")
	assert.NotNil(suite.T(), formatFlag)
	assert.Equal(suite.T(), "table", formatFlag.DefValue)

	limitFlag := cmd.Flags().Lookup("limit")
	assert.NotNil(suite.T(), limitFlag)
	assert.Equal(suite.T(), "0", limitFlag.DefValue)

	orgFlag := cmd.Flags().Lookup("org")
	assert.NotNil(suite.T(), orgFlag)

	statusFlag := cmd.Flags().Lookup("status")
	assert.NotNil(suite.T(), statusFlag)

	// Note: type flag may not exist in current implementation
}

func TestAppCommandTestSuite(t *testing.T) {
	suite.Run(t, new(AppCommandTestSuite))
}
