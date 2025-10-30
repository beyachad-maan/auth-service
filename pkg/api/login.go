package api

const ResourceTypeLogin = "Login"

type Login struct {
	Version      *string `json:"version,omitempty"`
	ResourceType *string `json:"resource_type,omitempty"`
	Username     string  `json:"username"`
	Password     string  `json:"password"`
}
