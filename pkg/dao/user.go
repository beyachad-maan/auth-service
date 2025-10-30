package dao

import (
	"context"
	"database/sql"
	"errors"

	"github.com/beyachad-maan/auth-service/pkg/models"
	"github.com/jmoiron/sqlx"
)

var ErrUserNotFound = errors.New("user not found")

type Users interface {
	CreateUser(ctx context.Context, user models.User) error
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	DeleteUserByID(ctx context.Context, id string) error
	GetUserByUserName(ctx context.Context, userName string) (*models.User, error)
	AddPointsToUserScoreById(ctx context.Context, id string, points int) error
}

type UsersSql struct {
	DB *sqlx.DB
}

func NewUsers(db *sqlx.DB) Users {
	return &UsersSql{
		DB: db,
	}
}

func (dao *UsersSql) CreateUser(ctx context.Context, user models.User) error {
	_, err := dao.DB.ExecContext(ctx, `INSERT INTO users (id, username, private_name, family_name, email, ethnicity, password, created_at) 
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		user.ID,
		user.Username,
		user.PrivateName,
		user.FamilyName,
		user.Email,
		user.Ethnicity,
		user.Password,
		user.CreatedAt,
	)
	return err
}

func (dao *UsersSql) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	row := dao.DB.QueryRowContext(ctx, "SELECT * FROM users WHERE id = $1", id)
	err := row.Scan(&user.ID,
		&user.Username,
		&user.PrivateName,
		&user.FamilyName,
		&user.Email,
		&user.Ethnicity,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (dao *UsersSql) GetUserByUserName(ctx context.Context, userName string) (*models.User, error) {
	var user models.User
	row := dao.DB.QueryRowContext(ctx, "SELECT * FROM users WHERE username = $1", userName)
	err := row.Scan(&user.ID,
		&user.Username,
		&user.PrivateName,
		&user.FamilyName,
		&user.Email,
		&user.Ethnicity,
		&user.Password,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (dao *UsersSql) DeleteUserByID(ctx context.Context, id string) error {
	_, err := dao.DB.ExecContext(ctx, "DELETE FROM users WHERE id = $1", id)
	return err
}

func (dao *UsersSql) AddPointsToUserScoreById(ctx context.Context, id string, points int) error {
	_, err := dao.DB.ExecContext(ctx, "UPDATE users SET score = score + $1 WHERE id = $2", points, id)
	return err
}
