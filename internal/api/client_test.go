package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"hawkop/internal/config"
)

// ClientTestSuite contains tests for the API client
type ClientTestSuite struct {
	suite.Suite
	client     *Client
	server     *httptest.Server
	testConfig *config.Config
}

// SetupSuite runs before all tests in the suite
func (suite *ClientTestSuite) SetupSuite() {
	// Create test config with mock credentials
	suite.testConfig = &config.Config{
		APIKey: "test-api-key",
		OrgID:  "test-org-id",
		JWT: &config.JWT{
			Token:     "test-jwt-token",
			ExpiresAt: time.Now().Add(1 * time.Hour),
		},
	}

	// Create test HTTP server
	suite.server = httptest.NewServer(http.HandlerFunc(suite.mockAPIHandler))

	// Create client with test server URL
	suite.client = NewClient(suite.testConfig)
	suite.client.SetBaseURL(suite.server.URL)
}

// TearDownSuite runs after all tests in the suite
func (suite *ClientTestSuite) TearDownSuite() {
	suite.server.Close()
}

// mockAPIHandler handles mock API responses for testing
func (suite *ClientTestSuite) mockAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.URL.Path {
	case "/api/v1/user":
		suite.handleMockUser(w, r)
	case "/api/v1/org/test-org-id/members":
		suite.handleMockMembers(w, r)
	case "/api/v1/org/test-org-id/teams":
		suite.handleMockTeams(w, r)
	case "/api/v2/org/test-org-id/apps":
		suite.handleMockApps(w, r)
	case "/api/v1/scan/test-org-id":
		suite.handleMockScans(w, r)
	case "/api/v1/auth/login":
		suite.handleMockAuth(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (suite *ClientTestSuite) handleMockAuth(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-ApiKey") != "test-api-key" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	auth := AuthResponse{
		Token:     "new-jwt-token",
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}
	json.NewEncoder(w).Encode(auth)
}

func (suite *ClientTestSuite) handleMockUser(w http.ResponseWriter, r *http.Request) {
	user := UserResponse{
		User: User{
			StackhawkId: "test-user-id",
			External: UserExternal{
				Email:    "test@example.com",
				FullName: "Test User",
				Organizations: []OrganizationMembership{
					{
						Organization: Organization{
							ID:   "test-org-id",
							Name: "Test Organization",
						},
						Role: "OWNER",
					},
				},
			},
		},
	}
	json.NewEncoder(w).Encode(user)
}

func (suite *ClientTestSuite) handleMockMembers(w http.ResponseWriter, r *http.Request) {
	// Check for pagination parameters
	pageSize := r.URL.Query().Get("pageSize")
	assert.Equal(suite.T(), "1000", pageSize, "Should use maximum page size")

	members := OrganizationMembersResponse{
		Users: []OrganizationMember{
			{
				StackhawkId: "user-1",
				External: &UserExternal{
					Email:    "user1@example.com",
					FullName: "User One",
					Organizations: []OrganizationMembership{
						{Role: "ADMIN"},
					},
				},
			},
			{
				StackhawkId: "user-2",
				External: &UserExternal{
					Email:    "user2@example.com",
					FullName: "User Two",
					Organizations: []OrganizationMembership{
						{Role: "MEMBER"},
					},
				},
			},
		},
	}
	json.NewEncoder(w).Encode(members)
}

func (suite *ClientTestSuite) handleMockTeams(w http.ResponseWriter, r *http.Request) {
	// Check for pagination parameters
	pageSize := r.URL.Query().Get("pageSize")
	assert.Equal(suite.T(), "1000", pageSize, "Should use maximum page size")

	teams := OrganizationTeamsResponse{
		Teams: []Team{
			{
				ID:   "team-1",
				Name: "Test Team 1",
				Users: []OrganizationMember{
					{StackhawkId: "user-1"},
				},
			},
			{
				ID:   "team-2",
				Name: "Test Team 2",
				Applications: []Application{
					{ID: "app-1", Name: "Test App"},
				},
			},
		},
	}
	json.NewEncoder(w).Encode(teams)
}

func (suite *ClientTestSuite) handleMockApps(w http.ResponseWriter, r *http.Request) {
	// Check for pagination parameters
	pageSize := r.URL.Query().Get("pageSize")
	assert.Equal(suite.T(), "1000", pageSize, "Should use maximum page size")

	apps := OrganizationApplicationsResponse{
		Applications: []AppApplication{
			{
				ApplicationID:     "app-1",
				Name:              "Test Application",
				ApplicationStatus: "ACTIVE",
				ApplicationType:   "STANDARD",
			},
		},
	}
	json.NewEncoder(w).Encode(apps)
}

func (suite *ClientTestSuite) handleMockScans(w http.ResponseWriter, r *http.Request) {
	// Check for pagination parameters
	pageSize := r.URL.Query().Get("pageSize")
	assert.Equal(suite.T(), "1000", pageSize, "Should use maximum page size")

	scans := OrganizationScansResponse{
		ApplicationScanResults: []ApplicationScanResult{
			{
				Scan: Scan{
					ID:              "scan-1",
					ApplicationID:   "app-1",
					ApplicationName: "Test App",
					Status:          "COMPLETED",
					Timestamp:       "1756596062834",
				},
				ScanDuration: "45",
				URLCount:     "10",
				AlertStats: &AlertStats{
					High:   2,
					Medium: 3,
					Low:    1,
					Total:  6,
				},
			},
		},
	}
	json.NewEncoder(w).Encode(scans)
}

// Test API client creation
func (suite *ClientTestSuite) TestNewClient() {
	client := NewClient(suite.testConfig)
	assert.NotNil(suite.T(), client)
	assert.Equal(suite.T(), DefaultBaseURL, client.BaseURL)
	assert.NotNil(suite.T(), client.HTTPClient)
}

// Test BuildStandardParams with defaults
func (suite *ClientTestSuite) TestBuildStandardParams_Defaults() {
	params := suite.client.BuildStandardParams(nil)

	assert.Equal(suite.T(), "1000", params["pageSize"])
	assert.Len(suite.T(), params, 1) // Only pageSize should be set
}

// Test BuildStandardParams with overrides
func (suite *ClientTestSuite) TestBuildStandardParams_Overrides() {
	overrides := map[string]string{
		"pageSize":  "500",
		"sortField": "timestamp",
		"sortDir":   "asc",
	}

	params := suite.client.BuildStandardParams(overrides)

	assert.Equal(suite.T(), "500", params["pageSize"])
	assert.Equal(suite.T(), "timestamp", params["sortField"])
	assert.Equal(suite.T(), "asc", params["sortDir"])
	assert.Len(suite.T(), params, 3)
}

// Test successful user retrieval
func (suite *ClientTestSuite) TestGetUser_Success() {
	user, err := suite.client.GetUser()

	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), user)
	assert.Equal(suite.T(), "test-user-id", user.StackhawkId)
	assert.Equal(suite.T(), "Test User", user.External.FullName)
}

// Test organization members listing
func (suite *ClientTestSuite) TestListOrganizationMembers_Success() {
	members, err := suite.client.ListOrganizationMembers("test-org-id")

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), members, 2)
	assert.Equal(suite.T(), "user1@example.com", members[0].External.Email)
	assert.Equal(suite.T(), "ADMIN", members[0].External.Organizations[0].Role)
}

// Test organization teams listing
func (suite *ClientTestSuite) TestListOrganizationTeams_Success() {
	teams, err := suite.client.ListOrganizationTeams("test-org-id")

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), teams, 2)
	assert.Equal(suite.T(), "Test Team 1", teams[0].Name)
	assert.Len(suite.T(), teams[0].Users, 1)
	assert.Len(suite.T(), teams[1].Applications, 1)
}

// Test organization applications listing
func (suite *ClientTestSuite) TestListOrganizationApplications_Success() {
	apps, err := suite.client.ListOrganizationApplications("test-org-id")

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), apps, 1)
	assert.Equal(suite.T(), "Test Application", apps[0].Name)
	assert.Equal(suite.T(), "ACTIVE", apps[0].ApplicationStatus)
}

// Test organization scans listing
func (suite *ClientTestSuite) TestListOrganizationScans_Success() {
	scans, err := suite.client.ListOrganizationScans("test-org-id")

	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), scans, 1)
	assert.Equal(suite.T(), "scan-1", scans[0].Scan.ID)
	assert.Equal(suite.T(), "COMPLETED", scans[0].Scan.Status)
	assert.Equal(suite.T(), 6, scans[0].AlertStats.Total)
}

// Test error handling for invalid organization
func (suite *ClientTestSuite) TestListOrganizationMembers_InvalidOrg() {
	_, err := suite.client.ListOrganizationMembers("invalid-org")

	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "not found (404)")
}

// Test rate limiting behavior
func (suite *ClientTestSuite) TestRateLimiting() {
	start := time.Now()

	// Make multiple requests
	_, _ = suite.client.GetUser()
	_, _ = suite.client.GetUser()
	_, _ = suite.client.GetUser()

	elapsed := time.Since(start)

	// Should take at least 334ms for 3 requests (167ms * 2 intervals)
	assert.GreaterOrEqual(suite.T(), elapsed, 334*time.Millisecond)
}

// Run the test suite
func TestClientTestSuite(t *testing.T) {
	suite.Run(t, new(ClientTestSuite))
}
