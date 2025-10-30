package outbound

import (
	"github.com/beyachad-maan/auth-service/pkg/api"
	"github.com/beyachad-maan/auth-service/pkg/models"
	"github.com/beyachad-maan/auth-service/pkg/ptr"
)

func MapUser(user models.User) api.User {
	return api.User{
		ResourceType: ptr.Addr(api.ResourceTypeUser),
		Version:      ptr.Addr(api.Version),
		ID:           ptr.Addr(user.ID.String()),
		Username:     user.Username,
		PrivateName:  user.PrivateName,
		FamilyName:   user.FamilyName,
		Email:        user.Email,
		Ethnicity:    user.Ethnicity,
	}
}
