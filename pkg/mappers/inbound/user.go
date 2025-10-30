package inbound

import (
	"github.com/beyachad-maan/auth-service/pkg/api"
	"github.com/beyachad-maan/auth-service/pkg/models"
)

func MapUser(user api.User) models.User {
	return models.User{
		Username:    user.Username,
		PrivateName: user.PrivateName,
		FamilyName:  user.FamilyName,
		Email:       user.Email,
		Ethnicity:   user.Ethnicity,
	}
}
