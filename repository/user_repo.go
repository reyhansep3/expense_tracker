package repository

import (
	"database/sql"
	"exp_tracker/models"
)

type UserRepository interface {
	Create(user *models.User) error
	IsUserExist(name, email string) (bool, error)
}

type UserRepo struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &UserRepo{db}
}

func (r *UserRepo) Create(b *models.User) error {
	query := `
		INSERT INTO users (id, name, email, password, token)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.Exec(query, b.ID, b.Name, b.Email, b.Password, b.Token)
	return err
}

func (r *UserRepo) IsUserExist(name, email string) (bool, error) {
	query := `SELECT EXISTS (
		SELECT 1 FROM users WHERE name = $1 OR email = $2
	)`
	var exists bool
	err := r.db.QueryRow(query, name, email).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
