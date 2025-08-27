package api

// No time import needed in this file

// Organization represents a StackHawk organization
type Organization struct {
	ID                string              `json:"id"`
	Name              string              `json:"name"`
	Plan              string              `json:"plan,omitempty"`
	CreatedTimestamp  string              `json:"createdTimestamp,omitempty"`
	Features          []string            `json:"features,omitempty"`
	Settings          *OrganizationSettings `json:"settings,omitempty"`
	Subscription      *Subscription       `json:"subscription,omitempty"`
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
	ID              string                    `json:"id"`
	Email           string                    `json:"email"`
	FirstName       string                    `json:"firstName"`
	LastName        string                    `json:"lastName"`
	FullName        string                    `json:"fullName"`
	AvatarUrl       string                    `json:"avatarUrl"`
	Organizations   []OrganizationMembership  `json:"organizations"`
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
	StackhawkId       string                     `json:"stackhawkId"`
	Provider          *Provider                  `json:"provider,omitempty"`
	External          *UserExternal              `json:"external,omitempty"`
	CreatedTimestamp  string                     `json:"createdTimestamp,omitempty"`
	Organization      *Organization              `json:"organization,omitempty"`
	Role              string                     `json:"role"`
	Features          []Feature                  `json:"features,omitempty"`
	Metadata          []Metadata                 `json:"metadata,omitempty"`
	Achievements      []Achievement              `json:"achievements,omitempty"`
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
	ID             string            `json:"id"`
	Name           string            `json:"name"`
	OrganizationID string            `json:"organizationId,omitempty"`
	Applications   []Application     `json:"applications,omitempty"`
	Users          []OrganizationMember `json:"users,omitempty"`
	CreatedTimestamp string          `json:"createdTimestamp,omitempty"`
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
	ApplicationID     string `json:"applicationId"`
	Name              string `json:"name"`
	Env               string `json:"env,omitempty"`
	EnvID             string `json:"envId,omitempty"`
	ApplicationStatus string `json:"applicationStatus,omitempty"`
	OrganizationID    string `json:"organizationId,omitempty"`
	ApplicationType   string `json:"applicationType,omitempty"`
	CloudScanTarget   interface{} `json:"cloudScanTarget,omitempty"`
}

// OrganizationApplicationsResponse represents the response from the /api/v2/org/{orgId}/apps endpoint
type OrganizationApplicationsResponse struct {
	Applications    []AppApplication `json:"applications,omitempty"`
	TotalCount      string           `json:"totalCount,omitempty"`
	CurrentPage     interface{}      `json:"currentPage,omitempty"`
	HasNext         bool             `json:"hasNext,omitempty"`
	NextPage        interface{}      `json:"nextPage,omitempty"`
	NextPageToken   string           `json:"nextPageToken,omitempty"`
	HasPrev         bool             `json:"hasPrev,omitempty"`
	PrevPage        interface{}      `json:"prevPage,omitempty"`
	PrevPageToken   string           `json:"prevPageToken,omitempty"`
}