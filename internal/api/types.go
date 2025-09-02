// Package api defines data structures for StackHawk API responses.
package api

// PaginationOptions represents pagination and sorting parameters
type PaginationOptions struct {
	PageSize  int    `json:"pageSize,omitempty"`
	PageToken string `json:"pageToken,omitempty"`
	Page      string `json:"page,omitempty"`
	SortField string `json:"sortField,omitempty"`
	SortDir   string `json:"sortDir,omitempty"`
}

// PaginationInfo represents pagination metadata in responses
type PaginationInfo struct {
	NextPageToken string      `json:"nextPageToken,omitempty"`
	PrevPageToken string      `json:"prevPageToken,omitempty"`
	TotalCount    string      `json:"totalCount,omitempty"`
	HasNext       bool        `json:"hasNext,omitempty"`
	HasPrev       bool        `json:"hasPrev,omitempty"`
	CurrentPage   interface{} `json:"currentPage,omitempty"`
	NextPage      interface{} `json:"nextPage,omitempty"`
	PrevPage      interface{} `json:"prevPage,omitempty"`
}

// Organization represents a StackHawk organization
type Organization struct {
	ID               string                `json:"id"`
	Name             string                `json:"name"`
	Plan             string                `json:"plan,omitempty"`
	CreatedTimestamp string                `json:"createdTimestamp,omitempty"`
	Features         []string              `json:"features,omitempty"`
	Settings         *OrganizationSettings `json:"settings,omitempty"`
	Subscription     *Subscription         `json:"subscription,omitempty"`
}

// OrganizationSettings represents organization configuration settings
type OrganizationSettings struct {
	// Add settings fields as needed based on API response
}

// Subscription represents billing/subscription information
type Subscription struct {
	Status string `json:"status,omitempty"`
	// Add other subscription fields as needed
}

// UserExternal represents external user info from providers
type UserExternal struct {
	ID            string                   `json:"id"`
	Email         string                   `json:"email"`
	FirstName     string                   `json:"firstName"`
	LastName      string                   `json:"lastName"`
	FullName      string                   `json:"fullName"`
	AvatarUrl     string                   `json:"avatarUrl"`
	Organizations []OrganizationMembership `json:"organizations"`
}

// OrganizationMembership represents a user's membership in an organization
type OrganizationMembership struct {
	Organization Organization `json:"organization"`
	Role         string       `json:"role"`
}

// User represents a StackHawk user with organization membership
type User struct {
	StackhawkId string       `json:"stackhawkId"`
	External    UserExternal `json:"external"`
}

// UserResponse represents the response from the /api/v1/user endpoint
type UserResponse struct {
	User User `json:"user"`
}

// OrganizationMember represents a user member of an organization
type OrganizationMember struct {
	StackhawkId      string        `json:"stackhawkId"`
	Provider         *Provider     `json:"provider,omitempty"`
	External         *UserExternal `json:"external,omitempty"`
	CreatedTimestamp string        `json:"createdTimestamp,omitempty"`
	Organization     *Organization `json:"organization,omitempty"`
	Role             string        `json:"role"`
	Features         []Feature     `json:"features,omitempty"`
	Metadata         []Metadata    `json:"metadata,omitempty"`
	Achievements     []Achievement `json:"achievements,omitempty"`
}

// Provider represents authentication provider information
type Provider struct {
	Slug     string `json:"slug"`
	ClientId string `json:"clientId"`
	Created  string `json:"created"`
}

// Feature represents access features for a user
type Feature struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

// Metadata represents organizational user metadata
type Metadata struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// Achievement represents product journey events
type Achievement struct {
	Achievement string `json:"achievement"`
	Timestamp   string `json:"timestamp"`
}

// OrganizationMembersResponse represents the response from the /api/v1/org/{orgId}/members endpoint
type OrganizationMembersResponse struct {
	Users         []OrganizationMember `json:"users,omitempty"`
	NextPageToken string               `json:"nextPageToken,omitempty"`
	TotalCount    string               `json:"totalCount,omitempty"`
}

// Team represents a StackHawk team within an organization
type Team struct {
	ID               string               `json:"id"`
	Name             string               `json:"name"`
	OrganizationID   string               `json:"organizationId,omitempty"`
	Applications     []Application        `json:"applications,omitempty"`
	Users            []OrganizationMember `json:"users,omitempty"`
	CreatedTimestamp string               `json:"createdTimestamp,omitempty"`
}

// Application represents a basic application reference in teams
type Application struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

// OrganizationTeamsResponse represents the response from the /api/v1/org/{orgId}/teams endpoint
type OrganizationTeamsResponse struct {
	Teams         []Team `json:"teams,omitempty"`
	NextPageToken string `json:"nextPageToken,omitempty"`
	TotalCount    string `json:"totalCount,omitempty"`
}

// AppApplication represents a StackHawk application
type AppApplication struct {
	ApplicationID     string      `json:"applicationId"`
	Name              string      `json:"name"`
	Env               string      `json:"env,omitempty"`
	EnvID             string      `json:"envId,omitempty"`
	ApplicationStatus string      `json:"applicationStatus,omitempty"`
	OrganizationID    string      `json:"organizationId,omitempty"`
	ApplicationType   string      `json:"applicationType,omitempty"`
	CloudScanTarget   interface{} `json:"cloudScanTarget,omitempty"`
}

// OrganizationApplicationsResponse represents the response from the /api/v2/org/{orgId}/apps endpoint
type OrganizationApplicationsResponse struct {
	Applications  []AppApplication `json:"applications,omitempty"`
	TotalCount    string           `json:"totalCount,omitempty"`
	CurrentPage   interface{}      `json:"currentPage,omitempty"`
	HasNext       bool             `json:"hasNext,omitempty"`
	NextPage      interface{}      `json:"nextPage,omitempty"`
	NextPageToken string           `json:"nextPageToken,omitempty"`
	HasPrev       bool             `json:"hasPrev,omitempty"`
	PrevPage      interface{}      `json:"prevPage,omitempty"`
	PrevPageToken string           `json:"prevPageToken,omitempty"`
}

// Scan represents a StackHawk scan
type Scan struct {
	ID              string `json:"id"`
	ApplicationID   string `json:"applicationId"`
	ApplicationName string `json:"applicationName"`
	Env             string `json:"env,omitempty"`
	Status          string `json:"status"`
	Timestamp       string `json:"timestamp"`
}

// ApplicationScanResult represents a scan result with metadata
type ApplicationScanResult struct {
	Scan         Scan        `json:"scan"`
	ScanDuration interface{} `json:"scanDuration,omitempty"`
	URLCount     interface{} `json:"urlCount,omitempty"`
	AlertStats   *AlertStats `json:"alertStats,omitempty"`
	AppHost      string      `json:"appHost,omitempty"`
	Timestamp    string      `json:"timestamp,omitempty"`
	PolicyName   string      `json:"policyName,omitempty"`
	Tags         interface{} `json:"tags,omitempty"`
	Metadata     interface{} `json:"metadata,omitempty"`
}

// AlertStats represents alert statistics for a scan
type AlertStats struct {
	High   int `json:"high,omitempty"`
	Medium int `json:"medium,omitempty"`
	Low    int `json:"low,omitempty"`
	Info   int `json:"info,omitempty"`
	Total  int `json:"total,omitempty"`
}

// OrganizationScansResponse represents the response from the /api/v1/scan/{orgId} endpoint
type OrganizationScansResponse struct {
	ApplicationScanResults []ApplicationScanResult `json:"applicationScanResults,omitempty"`
	NextPageToken          string                  `json:"nextPageToken,omitempty"`
	TotalCount             string                  `json:"totalCount,omitempty"`
}

// ScanAlert represents an alert/finding type in a scan
type ScanAlert struct {
	PluginID    string   `json:"pluginId"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Severity    string   `json:"severity"`
	References  []string `json:"references,omitempty"`
	URICount    int      `json:"uriCount,omitempty"`
	CWEID       string   `json:"cweId,omitempty"`
}

// ScanAlertsResponse represents the response from the /api/v1/scan/{scanId}/alerts endpoint
type ScanAlertsResponse struct {
	ApplicationScanResults []struct {
		ApplicationAlerts []ScanAlert `json:"applicationAlerts,omitempty"`
	} `json:"applicationScanResults,omitempty"`
	NextPageToken string `json:"nextPageToken,omitempty"`
}

// ScanAlertFinding represents a specific finding instance
type ScanAlertFinding struct {
	PluginID      string `json:"pluginId"`
	URI           string `json:"uri"`
	RequestMethod string `json:"requestMethod"`
	Status        string `json:"status"`
	MsgID         string `json:"msgId"`
}

// ScanAlertFindingsResponse represents the response from the /api/v1/scan/{scanId}/alert/{pluginId} endpoint
type ScanAlertFindingsResponse struct {
	Alert                    ScanAlert          `json:"alert,omitempty"`
	Category                 string             `json:"category,omitempty"`
	ApplicationScanAlertUris []ScanAlertFinding `json:"applicationScanAlertUris,omitempty"`
	AppHost                  string             `json:"appHost,omitempty"`
	TotalCount               string             `json:"totalCount,omitempty"`
	NextPageToken            string             `json:"nextPageToken,omitempty"`
}

// ScanMessage represents request/response data for a specific finding
type ScanMessage struct {
	ID             string `json:"id"`
	RequestHeader  string `json:"requestHeader,omitempty"`
	RequestBody    string `json:"requestBody,omitempty"`
	ResponseHeader string `json:"responseHeader,omitempty"`
	ResponseBody   string `json:"responseBody,omitempty"`
}

// ScanMessageResponse represents the response from the /api/v1/scan/{scanId}/uri/{alertUriId}/messages/{messageId} endpoint
type ScanMessageResponse struct {
	ScanMessage ScanMessage `json:"scanMessage,omitempty"`
	URI         string      `json:"uri,omitempty"`
	Evidence    string      `json:"evidence,omitempty"`
	Description string      `json:"description,omitempty"`
	Param       string      `json:"param,omitempty"`
}
