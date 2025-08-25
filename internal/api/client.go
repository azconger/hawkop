package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"hawkop/internal/config"
)

const (
	DefaultBaseURL = "https://api.stackhawk.com"
	AuthEndpoint   = "/api/v1/auth"
)

// Client represents the StackHawk API client
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	config     *config.Config
}

// AuthResponse represents the response from the authentication endpoint
type AuthResponse struct {
	AccessToken string    `json:"access_token"`
	ExpiresAt   time.Time `json:"expires_at"`
	TokenType   string    `json:"token_type"`
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

	// Prepare auth request body
	authReq := map[string]string{
		"api_key": c.config.APIKey,
	}

	reqBody, err := json.Marshal(authReq)
	if err != nil {
		return fmt.Errorf("failed to marshal auth request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", authURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create auth request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "hawkop-cli")

	// Make the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to authenticate: %w", err)
	}
	defer resp.Body.Close()

	// Check for success status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("authentication failed: HTTP %d", resp.StatusCode)
	}

	// Parse response
	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		return fmt.Errorf("failed to parse auth response: %w", err)
	}

	// Update JWT in config
	c.config.SetJWT(authResp.AccessToken, authResp.ExpiresAt)

	// Save config with new JWT
	if err := c.config.Save(); err != nil {
		return fmt.Errorf("failed to save JWT token: %w", err)
	}

	return nil
}

// DoAuthenticatedRequest performs an HTTP request with automatic JWT handling
func (c *Client) DoAuthenticatedRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	// Ensure we have a valid JWT
	if err := c.EnsureValidJWT(); err != nil {
		return nil, err
	}

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

	// Create request
	url := c.BaseURL + endpoint
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Authorization", "Bearer "+c.config.JWT.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "hawkop-cli")

	// Make the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Check if we got an unauthorized response (token might have expired)
	if resp.StatusCode == http.StatusUnauthorized {
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
	}

	return resp, nil
}

// Get performs a GET request with authentication
func (c *Client) Get(endpoint string) (*http.Response, error) {
	return c.DoAuthenticatedRequest("GET", endpoint, nil)
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