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

type UserCommandTestSuite struct {
	suite.Suite
	mockClient *api.MockClient
	origOutput io.Writer
	buffer     *bytes.Buffer
}

func (suite *UserCommandTestSuite) SetupTest() {
	suite.mockClient = api.NewMockClient()

	// Capture output
	suite.buffer = new(bytes.Buffer)
	suite.origOutput = os.Stdout
}

func (suite *UserCommandTestSuite) TearDownTest() {
	// Reset output if needed
}

func (suite *UserCommandTestSuite) TestUserListCommand_Success() {
	// Mock successful API response
	mockUsers := []api.OrganizationMember{
		{
			StackhawkId: "user-1",
			External: &api.UserExternal{
				Email:    "user1@example.com",
				FullName: "Test User 1",
				Organizations: []api.OrganizationMembership{
					{Role: "ADMIN"},
				},
			},
		},
	}

	suite.mockClient.On("ListOrganizationMembers", "test-org-id").Return(mockUsers, nil)

	cmd := userListCmd
	assert.Equal(suite.T(), "list", cmd.Use)
	assert.Contains(suite.T(), cmd.Short, "List users")
}

func (suite *UserCommandTestSuite) TestUserCommand_Structure() {
	// Test main user command structure
	assert.Equal(suite.T(), "user", userCmd.Use)
	assert.Contains(suite.T(), userCmd.Short, "Manage user")

	// Verify subcommands are registered
	subcommands := []string{}
	for _, cmd := range userCmd.Commands() {
		subcommands = append(subcommands, cmd.Use)
	}

	assert.Contains(suite.T(), subcommands, "list")
}

func (suite *UserCommandTestSuite) TestUserListFlags() {
	cmd := userListCmd

	// Test that all expected flags are present
	formatFlag := cmd.Flags().Lookup("format")
	assert.NotNil(suite.T(), formatFlag)
	assert.Equal(suite.T(), "table", formatFlag.DefValue)

	limitFlag := cmd.Flags().Lookup("limit")
	assert.NotNil(suite.T(), limitFlag)
	assert.Equal(suite.T(), "0", limitFlag.DefValue)

	orgFlag := cmd.Flags().Lookup("org")
	assert.NotNil(suite.T(), orgFlag)

	roleFlag := cmd.Flags().Lookup("role")
	assert.NotNil(suite.T(), roleFlag)
}

func TestUserCommandTestSuite(t *testing.T) {
	suite.Run(t, new(UserCommandTestSuite))
}
