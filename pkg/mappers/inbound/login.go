package inbound

import (
	"github.com/beyachad-maan/auth-service/pkg/api"
	"github.com/beyachad-maan/auth-service/pkg/models"
)

func MapLogin(login api.Login) models.Login {
	return models.Login{
		Username: login.Username,
		Password: login.Password,
	}
}
