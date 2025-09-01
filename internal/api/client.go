package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"hawkop/internal/config"
)

const (
	DefaultBaseURL = "https://api.stackhawk.com"
	AuthEndpoint   = "/api/v1/auth/login"

	// Pagination constants - use max page size to minimize API requests
	DefaultPageSize = 1000 // Use maximum to reduce API calls
	MaxPageSize     = 1000

	// Rate limiting constants
	MaxRequestsPerMinute = 360
	RetryAfterDefault    = 60 * time.Second
)

// Client represents the StackHawk API client
type Client struct {
	BaseURL     string
	HTTPClient  *http.Client
	config      *config.Config
	lastRequest time.Time
}

// AuthResponse represents the response from the authentication endpoint
type AuthResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at,omitempty"`
	TokenType string    `json:"token_type,omitempty"`
}

// NewClient creates a new StackHawk API client
func NewClient(cfg *config.Config) *Client {
	return &Client{
		BaseURL: DefaultBaseURL,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		config: cfg,
	}
}

// SetBaseURL updates the base URL for the API client
func (c *Client) SetBaseURL(baseURL string) {
	c.BaseURL = baseURL
}

// EnsureValidJWT checks if we have a valid JWT token and refreshes it if needed
func (c *Client) EnsureValidJWT() error {
	// Check if we need to refresh the JWT
	if !c.config.NeedsJWTRefresh() {
		return nil
	}

	// Check if we have valid credentials for authentication
	if !c.config.HasValidCredentials() {
		return fmt.Errorf("no API key configured - run 'hawkop init' to set up credentials")
	}

	// Authenticate to get a new JWT
	return c.authenticate()
}

// authenticate performs authentication with the StackHawk API to get a JWT token
func (c *Client) authenticate() error {
	authURL := c.BaseURL + AuthEndpoint

	// Create HTTP GET request with API key in X-ApiKey header (as per curl example)
	req, err := http.NewRequest("GET", authURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create auth request: %w", err)
	}

	req.Header.Set("X-ApiKey", c.config.APIKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "hawkop-cli")

	// Make the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}
	defer resp.Body.Close()

	// Check for success status
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("authentication failed: HTTP %d - %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse response
	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return fmt.Errorf("failed to parse auth response: %w", err)
	}

	// If no expiration is provided, set it to 30 minutes from now (as mentioned in the docs)
	expiresAt := authResp.ExpiresAt
	if expiresAt.IsZero() {
		expiresAt = time.Now().Add(30 * time.Minute)
	}

	// Update JWT in config
	c.config.SetJWT(authResp.Token, expiresAt)

	// Save config with new JWT
	if err := c.config.Save(); err != nil {
		return fmt.Errorf("failed to save JWT token: %w", err)
	}

	return nil
}

// DoAuthenticatedRequest performs an HTTP request with automatic JWT handling, rate limiting, and retry logic
func (c *Client) DoAuthenticatedRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	return c.DoAuthenticatedRequestWithParams(method, endpoint, body, nil)
}

// DoAuthenticatedRequestWithParams performs an HTTP request with pagination and query parameters
func (c *Client) DoAuthenticatedRequestWithParams(method, endpoint string, body interface{}, params map[string]string) (*http.Response, error) {
	// Ensure we have a valid JWT
	if err := c.EnsureValidJWT(); err != nil {
		return nil, err
	}

	// Rate limiting: ensure we don't exceed 360 requests per minute
	c.respectRateLimit()

	// Prepare request body
	var reqBody *bytes.Buffer
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(bodyBytes)
	} else {
		reqBody = &bytes.Buffer{}
	}

	// Build URL with query parameters
	reqURL := c.BaseURL + endpoint
	if params != nil && len(params) > 0 {
		u, err := url.Parse(reqURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse URL: %w", err)
		}

		q := u.Query()
		for key, value := range params {
			if value != "" {
				q.Set(key, value)
			}
		}
		u.RawQuery = q.Encode()
		reqURL = u.String()
	}

	// Create request
	req, err := http.NewRequest(method, reqURL, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers with Bearer JWT token
	req.Header.Set("Authorization", "Bearer "+c.config.JWT.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "hawkop-cli")

	// Make the request with retry logic
	resp, err := c.makeRequestWithRetry(req)
	if err != nil {
		return nil, err
	}

	// Update last request time for rate limiting
	c.lastRequest = time.Now()

	return resp, nil
}

// respectRateLimit implements basic rate limiting to stay under 360 requests/minute
func (c *Client) respectRateLimit() {
	// Simple rate limiting: ensure at least 167ms between requests (360/min = 6/sec)
	minInterval := 167 * time.Millisecond
	if !c.lastRequest.IsZero() {
		elapsed := time.Since(c.lastRequest)
		if elapsed < minInterval {
			time.Sleep(minInterval - elapsed)
		}
	}
}

// makeRequestWithRetry executes an HTTP request with retry logic for rate limiting and auth errors
func (c *Client) makeRequestWithRetry(req *http.Request) (*http.Response, error) {
	// Make the initial request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Handle different HTTP status codes
	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
		return resp, nil

	case http.StatusUnauthorized:
		resp.Body.Close()

		// Clear the JWT and try once more
		c.config.ClearJWT()
		if err := c.EnsureValidJWT(); err != nil {
			return nil, fmt.Errorf("failed to refresh token after 401: %w", err)
		}

		// Retry the request with new token
		req.Header.Set("Authorization", "Bearer "+c.config.JWT.Token)
		resp, err = c.HTTPClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("retry request failed: %w", err)
		}
		return resp, nil

	case http.StatusTooManyRequests:
		resp.Body.Close()

		// Check for Retry-After header
		retryAfter := RetryAfterDefault
		if retryHeader := resp.Header.Get("Retry-After"); retryHeader != "" {
			if seconds, err := strconv.Atoi(retryHeader); err == nil {
				retryAfter = time.Duration(seconds) * time.Second
			}
		}

		// Wait and retry once
		time.Sleep(retryAfter)
		resp, err = c.HTTPClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("retry after rate limit failed: %w", err)
		}
		return resp, nil

	case http.StatusBadRequest:
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("bad request (400): %s", string(bodyBytes))

	case http.StatusForbidden:
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("forbidden (403): insufficient permissions - %s", string(bodyBytes))

	case http.StatusNotFound:
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("not found (404): resource does not exist - %s", string(bodyBytes))

	case http.StatusConflict:
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("conflict (409): resource cannot be modified - %s", string(bodyBytes))

	case http.StatusUnprocessableEntity:
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("unprocessable entity (422): invalid input - %s", string(bodyBytes))

	default:
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("API error: HTTP %d - %s", resp.StatusCode, string(bodyBytes))
	}
}

// Get performs a GET request with authentication
func (c *Client) Get(endpoint string) (*http.Response, error) {
	return c.DoAuthenticatedRequest("GET", endpoint, nil)
}

// GetWithParams performs a GET request with authentication and query parameters
func (c *Client) GetWithParams(endpoint string, params map[string]string) (*http.Response, error) {
	return c.DoAuthenticatedRequestWithParams("GET", endpoint, nil, params)
}

// Post performs a POST request with authentication
func (c *Client) Post(endpoint string, body interface{}) (*http.Response, error) {
	return c.DoAuthenticatedRequest("POST", endpoint, body)
}

// Put performs a PUT request with authentication
func (c *Client) Put(endpoint string, body interface{}) (*http.Response, error) {
	return c.DoAuthenticatedRequest("PUT", endpoint, body)
}

// Delete performs a DELETE request with authentication
func (c *Client) Delete(endpoint string) (*http.Response, error) {
	return c.DoAuthenticatedRequest("DELETE", endpoint, nil)
}

// GetUser retrieves the current user information including organizations
func (c *Client) GetUser() (*User, error) {
	resp, err := c.Get("/api/v1/user")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API error: HTTP %d", resp.StatusCode)
	}

	var userResp UserResponse
	if err := json.NewDecoder(resp.Body).Decode(&userResp); err != nil {
		return nil, fmt.Errorf("failed to parse user response: %w", err)
	}

	return &userResp.User, nil
}

// ListOrganizations retrieves all organizations the user belongs to
func (c *Client) ListOrganizations() ([]Organization, error) {
	user, err := c.GetUser()
	if err != nil {
		return nil, fmt.Errorf("failed to get organizations: %w", err)
	}

	// Extract organizations from membership info
	organizations := make([]Organization, 0, len(user.External.Organizations))
	for _, membership := range user.External.Organizations {
		organizations = append(organizations, membership.Organization)
	}

	return organizations, nil
}

// ListOrganizationMembers retrieves all users/members in the specified organization
func (c *Client) ListOrganizationMembers(orgID string) ([]OrganizationMember, error) {
	endpoint := fmt.Sprintf("/api/v1/org/%s/members", orgID)

	// Use standard parameters with optimal defaults
	params := c.BuildStandardParams(nil)

	resp, err := c.GetWithParams(endpoint, params)
	if err != nil {
		return nil, err // Error handling now done in makeRequestWithRetry
	}
	defer resp.Body.Close()

	// Parse the wrapped response (users are in a "users" array)
	var wrappedResp OrganizationMembersResponse
	if err := json.NewDecoder(resp.Body).Decode(&wrappedResp); err != nil {
		return nil, fmt.Errorf("failed to parse organization members response: %w", err)
	}
	members := wrappedResp.Users
	return members, nil
}

// ListOrganizationTeams retrieves all teams in the specified organization
func (c *Client) ListOrganizationTeams(orgID string) ([]Team, error) {
	endpoint := fmt.Sprintf("/api/v1/org/%s/teams", orgID)

	// Use standard parameters with optimal defaults
	params := c.BuildStandardParams(nil)

	resp, err := c.GetWithParams(endpoint, params)
	if err != nil {
		return nil, err // Error handling now done in makeRequestWithRetry
	}
	defer resp.Body.Close()

	// Parse the response (teams are in a "teams" array)
	var teamsResp OrganizationTeamsResponse
	if err := json.NewDecoder(resp.Body).Decode(&teamsResp); err != nil {
		return nil, fmt.Errorf("failed to parse organization teams response: %w", err)
	}

	return teamsResp.Teams, nil
}

// ListOrganizationApplications retrieves all applications in the specified organization
func (c *Client) ListOrganizationApplications(orgID string) ([]AppApplication, error) {
	endpoint := fmt.Sprintf("/api/v2/org/%s/apps", orgID)

	// Use standard parameters with optimal defaults
	params := c.BuildStandardParams(nil)

	resp, err := c.GetWithParams(endpoint, params)
	if err != nil {
		return nil, err // Error handling now done in makeRequestWithRetry
	}
	defer resp.Body.Close()

	// Parse the response (applications are in an "applications" array)
	var appsResp OrganizationApplicationsResponse
	if err := json.NewDecoder(resp.Body).Decode(&appsResp); err != nil {
		return nil, fmt.Errorf("failed to parse organization applications response: %w", err)
	}

	return appsResp.Applications, nil
}

// ListOrganizationScans retrieves all scans for the specified organization
func (c *Client) ListOrganizationScans(orgID string) ([]ApplicationScanResult, error) {
	return c.ListOrganizationScansWithOptions(orgID, nil)
}

// ListOrganizationScansWithOptions retrieves scans with pagination and sorting options
func (c *Client) ListOrganizationScansWithOptions(orgID string, opts *PaginationOptions) ([]ApplicationScanResult, error) {
	endpoint := fmt.Sprintf("/api/v1/scan/%s", orgID)

	// Start with standard parameters (includes optimal pageSize=1000)
	overrides := make(map[string]string)

	// Apply pagination options as overrides
	if opts != nil {
		if opts.PageSize > 0 {
			if opts.PageSize > MaxPageSize {
				opts.PageSize = MaxPageSize
			}
			overrides["pageSize"] = strconv.Itoa(opts.PageSize)
		}
		if opts.PageToken != "" {
			overrides["pageToken"] = opts.PageToken
		}
		if opts.Page != "" {
			overrides["page"] = opts.Page
		}
		if opts.SortField != "" {
			overrides["sortField"] = opts.SortField
		}
		if opts.SortDir != "" {
			overrides["sortDir"] = opts.SortDir
		}
	}

	params := c.BuildStandardParams(overrides)

	resp, err := c.GetWithParams(endpoint, params)
	if err != nil {
		return nil, err // Error handling now done in makeRequestWithRetry
	}
	defer resp.Body.Close()

	// Parse the response
	var scansResp OrganizationScansResponse
	if err := json.NewDecoder(resp.Body).Decode(&scansResp); err != nil {
		return nil, fmt.Errorf("failed to parse organization scans response: %w", err)
	}

	return scansResp.ApplicationScanResults, nil
}

// GetScanAlerts retrieves alerts for a specific scan
func (c *Client) GetScanAlerts(scanID string) ([]ScanAlert, error) {
	endpoint := fmt.Sprintf("/api/v1/scan/%s/alerts", scanID)

	resp, err := c.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get scan alerts: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API error: HTTP %d - %s", resp.StatusCode, string(bodyBytes))
	}

	// Parse the response
	var alertsResp ScanAlertsResponse
	if err := json.NewDecoder(resp.Body).Decode(&alertsResp); err != nil {
		return nil, fmt.Errorf("failed to parse scan alerts response: %w", err)
	}

	// Extract alerts from nested structure
	var alerts []ScanAlert
	for _, result := range alertsResp.ApplicationScanResults {
		alerts = append(alerts, result.ApplicationAlerts...)
	}

	return alerts, nil
}

// BuildStandardParams creates optimized API parameters with smart defaults
func (c *Client) BuildStandardParams(overrides map[string]string) map[string]string {
	params := map[string]string{
		"pageSize": strconv.Itoa(DefaultPageSize), // Always 1000 for efficiency
	}

	// Apply overrides for special cases
	for key, value := range overrides {
		if value != "" {
			params[key] = value
		}
	}

	return params
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
