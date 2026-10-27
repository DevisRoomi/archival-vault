package models

import "time"

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role      string    `json:"role"` // ARCHIVIST, OPERATOR, VIEWER
	CreatedAt time.Time `json:"created_at"`
}

type Record struct {
	ID             int64     `json:"id"`
	Code           string    `json:"code"`           // e.g. ARC-2026-001
	Title          string    `json:"title"`
	Classification string    `json:"classification"` // PUBLIC, INTERNAL, RESTRICTED, CONFIDENTIAL
	Category       string    `json:"category"`       // SPEC, INCIDENT, SCHEMA, DISPATCH
	Status         string    `json:"status"`         // ACTIVE, ARCHIVED, DRAFT
	Tags           string    `json:"tags"`           // comma-separated tags
	Summary        string    `json:"summary"`
	Payload        string    `json:"payload"`        // raw text / markdown / JSON payload
	CreatedBy      string    `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type RecordRequest struct {
	Title          string `json:"title"`
	Classification string `json:"classification"`
	Category       string `json:"category"`
	Status         string `json:"status"`
	Tags           string `json:"tags"`
	Summary        string `json:"summary"`
	Payload        string `json:"payload"`
}

type StatsResponse struct {
	TotalRecords    int            `json:"total_records"`
	ByClass         map[string]int `json:"by_classification"`
	ByCategory      map[string]int `json:"by_category"`
	ByStatus        map[string]int `json:"by_status"`
	LatestCode      string         `json:"latest_code"`
	UptimeSeconds   int64          `json:"uptime_seconds"`
}

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}
