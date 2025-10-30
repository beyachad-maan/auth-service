package api

import (
	"time"
)

const ResourceTypeUser = "User"

type User struct {
	Version      *string    `json:"version,omitempty"`
	ResourceType *string    `json:"resource_type,omitempty"`
	ID           *string    `json:"id,omitempty"`
	Username     string     `json:"username"`
	PrivateName  string     `json:"private_name"`
	FamilyName   string     `json:"family_name"`
	Email        string     `json:"email"`
	Ethnicity    string     `json:"ethnicity"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
}
