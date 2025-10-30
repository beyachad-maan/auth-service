package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `db:"id"`
	Username    string    `db:"username"`
	PrivateName string    `db:"private_name"`
	FamilyName  string    `db:"family_name"`
	Email       string    `db:"email"`
	Ethnicity   string    `db:"ethnicity"`
	Password    string    `db:"password"`
	CreatedAt   time.Time `db:"created_at"`
}

func NewID() uuid.UUID {
	return uuid.New()
}
