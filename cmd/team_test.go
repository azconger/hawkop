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

type TeamCommandTestSuite struct {
	suite.Suite
	mockClient *api.MockClient
	origOutput io.Writer
	buffer     *bytes.Buffer
}

func (suite *TeamCommandTestSuite) SetupTest() {
	suite.mockClient = api.NewMockClient()

	// Capture output
	suite.buffer = new(bytes.Buffer)
	suite.origOutput = os.Stdout
}

func (suite *TeamCommandTestSuite) TearDownTest() {
	// Reset output if needed
}

func (suite *TeamCommandTestSuite) TestTeamListCommand_Success() {
	// Mock successful API response
	mockTeams := []api.Team{
		{
			ID:   "team-1",
			Name: "Test Team",
			Users: []api.OrganizationMember{
				{StackhawkId: "user-1"},
			},
		},
	}

	suite.mockClient.On("ListOrganizationTeams", "test-org-id").Return(mockTeams, nil)

	cmd := teamListCmd
	assert.Equal(suite.T(), "list", cmd.Use)
	assert.Contains(suite.T(), cmd.Short, "List teams")
}

func (suite *TeamCommandTestSuite) TestTeamCommand_Structure() {
	// Test main team command structure
	assert.Equal(suite.T(), "team", teamCmd.Use)
	assert.Contains(suite.T(), teamCmd.Short, "Manage team")

	// Verify subcommands are registered
	subcommands := []string{}
	for _, cmd := range teamCmd.Commands() {
		subcommands = append(subcommands, cmd.Use)
	}

	assert.Contains(suite.T(), subcommands, "list")
}

func (suite *TeamCommandTestSuite) TestTeamListFlags() {
	cmd := teamListCmd

	// Test that all expected flags are present
	formatFlag := cmd.Flags().Lookup("format")
	assert.NotNil(suite.T(), formatFlag)
	assert.Equal(suite.T(), "table", formatFlag.DefValue)

	limitFlag := cmd.Flags().Lookup("limit")
	assert.NotNil(suite.T(), limitFlag)
	assert.Equal(suite.T(), "0", limitFlag.DefValue)

	orgFlag := cmd.Flags().Lookup("org")
	assert.NotNil(suite.T(), orgFlag)
}

func TestTeamCommandTestSuite(t *testing.T) {
	suite.Run(t, new(TeamCommandTestSuite))
}
