package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/stretchr/testify/mock"
)

// MockClient implements the Client interface for testing
type MockClient struct {
	mock.Mock
	BaseURL string
}

// NewMockClient creates a new mock client
func NewMockClient() *MockClient {
	return &MockClient{
		BaseURL: "http://localhost:8080",
	}
}

// SetBaseURL sets the base URL for the mock client
func (m *MockClient) SetBaseURL(url string) {
	m.BaseURL = url
}

// GetUser mocks the GetUser method
func (m *MockClient) GetUser() (*User, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

// ListOrganizations mocks the ListOrganizations method
func (m *MockClient) ListOrganizations() ([]Organization, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Organization), args.Error(1)
}

// ListOrganizationMembers mocks the ListOrganizationMembers method
func (m *MockClient) ListOrganizationMembers(orgID string) ([]OrganizationMember, error) {
	args := m.Called(orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]OrganizationMember), args.Error(1)
}

// ListOrganizationTeams mocks the ListOrganizationTeams method
func (m *MockClient) ListOrganizationTeams(orgID string) ([]Team, error) {
	args := m.Called(orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]Team), args.Error(1)
}

// ListOrganizationApplications mocks the ListOrganizationApplications method
func (m *MockClient) ListOrganizationApplications(orgID string) ([]AppApplication, error) {
	args := m.Called(orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]AppApplication), args.Error(1)
}

// ListOrganizationScans mocks the ListOrganizationScans method
func (m *MockClient) ListOrganizationScans(orgID string) ([]ApplicationScanResult, error) {
	args := m.Called(orgID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]ApplicationScanResult), args.Error(1)
}

// GetScanAlerts mocks the GetScanAlerts method
func (m *MockClient) GetScanAlerts(scanID string) ([]ScanAlert, error) {
	args := m.Called(scanID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]ScanAlert), args.Error(1)
}

// MockAPIServer provides a test HTTP server with mock responses
type MockAPIServer struct {
	Server *httptest.Server
}

// NewMockAPIServer creates a new mock API server with predefined responses
func NewMockAPIServer() *MockAPIServer {
	server := httptest.NewServer(http.HandlerFunc(mockAPIHandler))
	return &MockAPIServer{
		Server: server,
	}
}

// Close shuts down the mock server
func (m *MockAPIServer) Close() {
	m.Server.Close()
}

// URL returns the mock server URL
func (m *MockAPIServer) URL() string {
	return m.Server.URL
}

// mockAPIHandler handles mock API responses
func mockAPIHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.URL.Path {
	case "/api/v1/user":
		handleMockUser(w, r)
	case "/api/v1/orgs":
		handleMockOrganizations(w, r)
	case "/api/v1/org/test-org-id/members":
		handleMockMembers(w, r)
	case "/api/v1/org/test-org-id/teams":
		handleMockTeams(w, r)
	case "/api/v2/org/test-org-id/apps":
		handleMockApps(w, r)
	case "/api/v1/scan/test-org-id":
		handleMockScans(w, r)
	case "/api/v1/auth/login":
		handleMockAuth(w, r)
	default:
		http.NotFound(w, r)
	}
}

func handleMockAuth(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("X-ApiKey") == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	auth := AuthResponse{
		Token:     "mock-jwt-token",
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}
	json.NewEncoder(w).Encode(auth)
}

func handleMockUser(w http.ResponseWriter, r *http.Request) {
	user := UserResponse{
		User: User{
			StackhawkId: "mock-user-id",
			External: UserExternal{
				Email:    "mock@example.com",
				FullName: "Mock User",
				Organizations: []OrganizationMembership{
					{
						Organization: Organization{
							ID:   "test-org-id",
							Name: "Mock Organization",
						},
						Role: "OWNER",
					},
				},
			},
		},
	}
	json.NewEncoder(w).Encode(user)
}

func handleMockOrganizations(w http.ResponseWriter, r *http.Request) {
	orgs := []Organization{
		{
			ID:   "org-1",
			Name: "Mock Organization 1",
		},
		{
			ID:   "org-2",
			Name: "Mock Organization 2",
		},
	}
	json.NewEncoder(w).Encode(orgs)
}

func handleMockMembers(w http.ResponseWriter, r *http.Request) {
	members := OrganizationMembersResponse{
		Users: []OrganizationMember{
			{
				StackhawkId: "user-1",
				External: &UserExternal{
					Email:    "user1@mock.com",
					FullName: "Mock User 1",
					Organizations: []OrganizationMembership{
						{Role: "ADMIN"},
					},
				},
			},
			{
				StackhawkId: "user-2",
				External: &UserExternal{
					Email:    "user2@mock.com",
					FullName: "Mock User 2",
					Organizations: []OrganizationMembership{
						{Role: "MEMBER"},
					},
				},
			},
		},
	}
	json.NewEncoder(w).Encode(members)
}

func handleMockTeams(w http.ResponseWriter, r *http.Request) {
	teams := OrganizationTeamsResponse{
		Teams: []Team{
			{
				ID:   "team-1",
				Name: "Mock Team 1",
				Users: []OrganizationMember{
					{StackhawkId: "user-1"},
				},
			},
			{
				ID:   "team-2",
				Name: "Mock Team 2",
				Applications: []Application{
					{ID: "app-1", Name: "Mock App"},
				},
			},
		},
	}
	json.NewEncoder(w).Encode(teams)
}

func handleMockApps(w http.ResponseWriter, r *http.Request) {
	apps := OrganizationApplicationsResponse{
		Applications: []AppApplication{
			{
				ApplicationID:     "app-1",
				Name:              "Mock Application",
				ApplicationStatus: "ACTIVE",
				ApplicationType:   "STANDARD",
			},
		},
	}
	json.NewEncoder(w).Encode(apps)
}

func handleMockScans(w http.ResponseWriter, r *http.Request) {
	scans := OrganizationScansResponse{
		ApplicationScanResults: []ApplicationScanResult{
			{
				Scan: Scan{
					ID:              "scan-1",
					ApplicationID:   "app-1",
					ApplicationName: "Mock App",
					Status:          "COMPLETED",
					Timestamp:       "1756596062834",
					Env:             "production",
				},
				ScanDuration: "45",
				URLCount:     "10",
				AlertStats: &AlertStats{
					High:   2,
					Medium: 3,
					Low:    1,
					Info:   0,
					Total:  6,
				},
			},
		},
	}
	json.NewEncoder(w).Encode(scans)
}
