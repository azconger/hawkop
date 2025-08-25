package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Config represents the hawkop configuration
type Config struct {
	APIKey string `json:"api_key,omitempty" yaml:"api_key,omitempty"`
	OrgID  string `json:"org_id,omitempty" yaml:"org_id,omitempty"`
	JWT    *JWT   `json:"jwt,omitempty" yaml:"jwt,omitempty"`
}

// JWT represents a JSON Web Token with expiration
type JWT struct {
	Token     string    `json:"token" yaml:"token"`
	ExpiresAt time.Time `json:"expires_at" yaml:"expires_at"`
}

// IsExpired checks if the JWT token has expired
func (j *JWT) IsExpired() bool {
	if j == nil {
		return true
	}
	return time.Now().After(j.ExpiresAt)
}

// IsValid checks if the JWT exists and is not expired
func (j *JWT) IsValid() bool {
	return j != nil && j.Token != "" && !j.IsExpired()
}

var (
	configDir  string
	configFile string
)

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(fmt.Sprintf("unable to get user home directory: %v", err))
	}
	
	configDir = filepath.Join(homeDir, ".config", "hawkop")
	configFile = filepath.Join(configDir, "config.json")
}

// GetConfigDir returns the configuration directory path
func GetConfigDir() string {
	return configDir
}

// GetConfigFile returns the configuration file path
func GetConfigFile() string {
	return configFile
}

// Load reads and parses the configuration file
func Load() (*Config, error) {
	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// Check if config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		// Return empty config if file doesn't exist
		return &Config{}, nil
	}

	// Read config file
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// Save writes the configuration to the config file
func (c *Config) Save() error {
	// Ensure config directory exists
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Marshal to JSON with indentation for readability
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file with restricted permissions
	if err := os.WriteFile(configFile, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// SetAPIKey updates the API key in the configuration
func (c *Config) SetAPIKey(apiKey string) {
	c.APIKey = apiKey
	// Clear JWT when API key changes
	c.JWT = nil
}

// SetOrgID updates the organization ID in the configuration
func (c *Config) SetOrgID(orgID string) {
	c.OrgID = orgID
}

// SetJWT updates the JWT token and expiration in the configuration
func (c *Config) SetJWT(token string, expiresAt time.Time) {
	c.JWT = &JWT{
		Token:     token,
		ExpiresAt: expiresAt,
	}
}

// ClearJWT removes the JWT token from the configuration
func (c *Config) ClearJWT() {
	c.JWT = nil
}

// HasValidCredentials checks if the config has required credentials for API access
func (c *Config) HasValidCredentials() bool {
	return c.APIKey != ""
}

// NeedsJWTRefresh checks if a new JWT token should be obtained
func (c *Config) NeedsJWTRefresh() bool {
	return c.HasValidCredentials() && (c.JWT == nil || c.JWT.IsExpired())
}