package outbound

import (
	"github.com/beyachad-maan/auth-service/pkg/api"
	"github.com/beyachad-maan/auth-service/pkg/models"
	"github.com/beyachad-maan/auth-service/pkg/ptr"
)

func MapLogin(login models.Login) api.Login {
	return api.Login{
		Version:      ptr.Addr[string](api.Version),
		ResourceType: ptr.Addr[string](api.ResourceTypeLogin),
		Username:     login.Username,
		Password:     login.Password,
	}
}
