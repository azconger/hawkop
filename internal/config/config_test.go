package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

func (suite *ConfigTestSuite) TestNewConfig() {
	cfg := &Config{
		APIKey: "test-api-key",
		OrgID:  "test-org-id",
		JWT: &JWT{
			Token:     "test-jwt-token",
			ExpiresAt: time.Now().Add(1 * time.Hour),
		},
	}

	assert.Equal(suite.T(), "test-api-key", cfg.APIKey)
	assert.Equal(suite.T(), "test-org-id", cfg.OrgID)
	assert.NotNil(suite.T(), cfg.JWT)
}

func (suite *ConfigTestSuite) TestJWT_IsExpired() {
	// Test expired JWT
	expiredJWT := &JWT{
		Token:     "expired-token",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	assert.True(suite.T(), expiredJWT.IsExpired())

	// Test valid JWT
	validJWT := &JWT{
		Token:     "valid-token",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	assert.False(suite.T(), validJWT.IsExpired())

	// Test nil JWT
	var nilJWT *JWT
	assert.True(suite.T(), nilJWT.IsExpired())
}

func (suite *ConfigTestSuite) TestConfig_NeedsJWTRefresh() {
	cfg := &Config{APIKey: "test-key"}

	// No JWT
	assert.True(suite.T(), cfg.NeedsJWTRefresh())

	// Expired JWT
	cfg.JWT = &JWT{
		Token:     "expired",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	assert.True(suite.T(), cfg.NeedsJWTRefresh())

	// Valid JWT
	cfg.JWT = &JWT{
		Token:     "valid",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	assert.False(suite.T(), cfg.NeedsJWTRefresh())
}

func (suite *ConfigTestSuite) TestConfig_HasValidCredentials() {
	cfg := &Config{}

	// No credentials
	assert.False(suite.T(), cfg.HasValidCredentials())

	// Only API key - this should be valid credentials
	cfg.APIKey = "test-key"
	assert.True(suite.T(), cfg.HasValidCredentials())

	// API key with expired JWT - still valid credentials (JWT state doesn't matter)
	cfg.JWT = &JWT{
		Token:     "expired",
		ExpiresAt: time.Now().Add(-1 * time.Hour),
	}
	assert.True(suite.T(), cfg.HasValidCredentials())

	// API key with valid JWT - still valid credentials
	cfg.JWT = &JWT{
		Token:     "valid",
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}
	assert.True(suite.T(), cfg.HasValidCredentials())
}

func (suite *ConfigTestSuite) TestSetAPIKey() {
	cfg := &Config{
		APIKey: "old-key",
		JWT: &JWT{
			Token:     "old-token",
			ExpiresAt: time.Now().Add(1 * time.Hour),
		},
	}

	cfg.SetAPIKey("new-key")

	assert.Equal(suite.T(), "new-key", cfg.APIKey)
	assert.Nil(suite.T(), cfg.JWT) // JWT should be cleared when API key changes
}

func (suite *ConfigTestSuite) TestOrgIDManagement() {
	cfg := &Config{}

	// Test setting org ID
	cfg.OrgID = "test-org-id"
	assert.Equal(suite.T(), "test-org-id", cfg.OrgID)

	// Test clearing org ID
	cfg.OrgID = ""
	assert.Empty(suite.T(), cfg.OrgID)
}

func (suite *ConfigTestSuite) TestGetConfigPaths() {
	// Test that config path functions return non-empty strings
	configDir := GetConfigDir()
	assert.NotEmpty(suite.T(), configDir)
	assert.Contains(suite.T(), configDir, "hawkop")

	configFile := GetConfigFile()
	assert.NotEmpty(suite.T(), configFile)
	assert.Contains(suite.T(), configFile, "config.yaml")
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
