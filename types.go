package sdk

// ---------------------------------------------------------------------------
// Auth
// ---------------------------------------------------------------------------

// AuthTokens is returned by authentication endpoints.
type AuthTokens struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

// ---------------------------------------------------------------------------
// Domain models
// ---------------------------------------------------------------------------

// User represents an authenticated user account.
type User struct {
	ID            string  `json:"id"`
	Email         string  `json:"email"`
	Name          string  `json:"name"`
	Plan          string  `json:"plan"`
	DataRegion    string  `json:"data_region"`
	Country       *string `json:"country"`
	TosAcceptedAt *string `json:"tos_accepted_at"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

// Roost is a webhook endpoint that captures incoming requests.
type Roost struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Secret      *string `json:"secret"`
	SourceType  *string `json:"source_type"`
	IsActive    bool    `json:"is_active"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// Pigeon is a captured webhook request.
type Pigeon struct {
	ID             string         `json:"id"`
	RoostID        string         `json:"roost_id"`
	SourceIP       string         `json:"source_ip"`
	RequestMethod  string         `json:"request_method"`
	ContentType    string         `json:"content_type"`
	Headers        map[string]any `json:"headers"`
	BodyJSON       any            `json:"body_json"`
	BodyRaw        []int          `json:"body_raw"`
	RequestQuery   map[string]any `json:"request_query"`
	ReplayedFrom   *string        `json:"replayed_from"`
	DeliveryStatus string         `json:"delivery_status"`
	ReceivedAt     string         `json:"received_at"`
}

// Destination is a forwarding target attached to a roost.
type Destination struct {
	ID               string         `json:"id"`
	RoostID          string         `json:"roost_id"`
	Name             string         `json:"name"`
	DestinationType  string         `json:"destination_type"`
	Config           map[string]any `json:"config"`
	FilterExpression string         `json:"filter_expression"`
	Template         string         `json:"template"`
	RetryMax         int            `json:"retry_max"`
	RetryDelayMs     int            `json:"retry_delay_ms"`
	RetryMultiplier  float64        `json:"retry_multiplier"`
	IsPaused         bool           `json:"is_paused"`
	IsVerified       bool           `json:"is_verified"`
	CreatedAt        string         `json:"created_at"`
	UpdatedAt        string         `json:"updated_at"`
}

// DeliveryAttempt is a single attempt to deliver a pigeon to a destination.
type DeliveryAttempt struct {
	ID              string            `json:"id"`
	PigeonID        string            `json:"pigeon_id"`
	DestinationID   string            `json:"destination_id"`
	Status          string            `json:"status"`
	AttemptNumber   int               `json:"attempt_number"`
	ResponseStatus  *int              `json:"response_status"`
	ResponseBody    *string           `json:"response_body"`
	ResponseHeaders map[string]string `json:"response_headers,omitempty"`
	ErrorMessage    *string           `json:"error_message"`
	AttemptedAt     string            `json:"attempted_at"`
	NextRetryAt     *string           `json:"next_retry_at"`
}

// ---------------------------------------------------------------------------
// API Keys
// ---------------------------------------------------------------------------

// ApiKeyResponse is an API key without the full key value.
type ApiKeyResponse struct {
	ID        string  `json:"id"`
	KeyPrefix string  `json:"key_prefix"`
	Name      string  `json:"name"`
	LastUsed  *string `json:"last_used"`
	RevokedAt *string `json:"revoked_at"`
	CreatedAt string  `json:"created_at"`
}

// ApiKeyCreatedResponse includes the full key (shown only once).
type ApiKeyCreatedResponse struct {
	ID        string `json:"id"`
	Key       string `json:"key"`
	KeyPrefix string `json:"key_prefix"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
}

// ---------------------------------------------------------------------------
// Mutation requests
// ---------------------------------------------------------------------------

// CreateRoost is the body for POST /v1/roosts.
type CreateRoost struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Secret      *string `json:"secret,omitempty"`
	SourceType  *string `json:"source_type,omitempty"`
}

// UpdateRoost is the body for PATCH /v1/roosts/:id.
type UpdateRoost struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Secret      *string `json:"secret,omitempty"`
	SourceType  *string `json:"source_type,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

// CreateDestination is the body for POST /v1/roosts/:id/destinations.
type CreateDestination struct {
	DestinationType  string         `json:"destination_type"`
	Name             *string        `json:"name,omitempty"`
	Config           map[string]any `json:"config,omitempty"`
	FilterExpression *string        `json:"filter_expression,omitempty"`
	Template         *string        `json:"template,omitempty"`
	RetryMax         *int           `json:"retry_max,omitempty"`
	RetryDelayMs     *int           `json:"retry_delay_ms,omitempty"`
	RetryMultiplier  *float64       `json:"retry_multiplier,omitempty"`
}

// UpdateDestination is the body for PATCH /v1/destinations/:id.
type UpdateDestination struct {
	Name             *string        `json:"name,omitempty"`
	Config           map[string]any `json:"config,omitempty"`
	FilterExpression *string        `json:"filter_expression,omitempty"`
	Template         *string        `json:"template,omitempty"`
	TransformType    *string        `json:"transform_type,omitempty"`
	TransformCode    *string        `json:"transform_code,omitempty"`
}

// PauseInput is the body for PATCH /v1/destinations/:id/pause.
type PauseInput struct {
	IsPaused bool `json:"is_paused"`
}

// PauseResponse is returned by PATCH /v1/destinations/:id/pause.
type PauseResponse struct {
	IsPaused bool `json:"is_paused"`
}

// UpdateProfileRequest is the body for PATCH /v1/me.
type UpdateProfileRequest struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}

// CreateApiKeyRequest is the body for POST /v1/api-keys.
type CreateApiKeyRequest struct {
	Name *string `json:"name,omitempty"`
}

// UpdateApiKeyRequest is the body for PATCH /v1/api-keys/:id.
type UpdateApiKeyRequest struct {
	Name string `json:"name"`
}

// ---------------------------------------------------------------------------
// Responses
// ---------------------------------------------------------------------------

// ReplayResponse is returned by POST /v1/pigeons/:id/replay.
type ReplayResponse struct {
	Replayed         bool   `json:"replayed"`
	PigeonID         string `json:"pigeon_id"`
	DeliveryAttempts int    `json:"delivery_attempts"`
}

// ---------------------------------------------------------------------------
// Pagination
// ---------------------------------------------------------------------------

// PaginatedPigeons is a paginated list of pigeons.
type PaginatedPigeons struct {
	Data       []Pigeon `json:"data"`
	NextCursor *string  `json:"next_cursor"`
	HasMore    bool     `json:"has_more"`
}

// PaginatedDeliveryAttempts is a paginated list of delivery attempts.
type PaginatedDeliveryAttempts struct {
	Data       []DeliveryAttempt `json:"data"`
	NextCursor *string           `json:"next_cursor"`
	HasMore    bool              `json:"has_more"`
}

// ---------------------------------------------------------------------------
// Templates
// ---------------------------------------------------------------------------

// Template is a reusable template for formatting webhook payloads.
type Template struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Body        string `json:"body"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreateTemplate is the body for POST /v1/templates.
type CreateTemplate struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
	Body        *string `json:"body,omitempty"`
}

// UpdateTemplate is the body for PATCH /v1/templates/:id.
type UpdateTemplate struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
	Body        *string `json:"body,omitempty"`
}

// PreviewTemplateRequest is the body for POST /v1/templates/preview.
type PreviewTemplateRequest struct {
	Body     string `json:"body"`
	PigeonID string `json:"pigeon_id"`
}

// PreviewTemplateResponse is returned by POST /v1/templates/preview.
type PreviewTemplateResponse struct {
	Rendered string `json:"rendered"`
}
