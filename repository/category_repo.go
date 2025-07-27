package repository

import (
	"database/sql"
	"exp_tracker/models"
	"fmt"
)

type CategoryRepository interface {
	Create(category *models.Categories) error
	GetAllData(userId int) ([]models.Categories, error)
	GetDataByID(id int, userId int) (*models.Categories, error)
	DeleteDataByID(id int, userId int) error
}

type CategoryRepo struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &CategoryRepo{db}
}

func (r *CategoryRepo) Create(b *models.Categories) error {
	query := `
		INSERT INTO categories (id, user_id, category_name, create_at, create_by, modified_at, modified_by)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(query, b.ID, b.UserId, b.CategoryName, b.CreateAt, b.CreateBy, b.ModifiedAt, b.ModifiedBy)
	if err != nil {
		return fmt.Errorf("failed to insert category: %w", err)
	}

	return nil
}

func (r *CategoryRepo) GetAllData(userId int) ([]models.Categories, error) {
	query := `
		SELECT * FROM categories WHERE user_id = $1
	`

	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get all categories: %w", err)
	}
	defer rows.Close()

	var categories []models.Categories
	for rows.Next() {
		var c models.Categories
		err := rows.Scan(&c.ID, &c.UserId, &c.CategoryName, &c.CreateAt, &c.CreateBy, &c.ModifiedAt, &c.ModifiedBy)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}

func (r *CategoryRepo) GetDataByID(id int, userId int) (*models.Categories, error) {
	query := `
		SELECT * FROM categories WHERE id = $1 AND user_id = $2
	`

	var c models.Categories

	err := r.db.QueryRow(query, id, userId).Scan(&c.ID, &c.UserId, &c.CategoryName, &c.CreateAt, &c.CreateBy, &c.ModifiedAt, &c.ModifiedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("category not found")
		}
		return nil, fmt.Errorf("failed to get category by id")
	}
	return &c, nil
}

func (r *CategoryRepo) DeleteDataByID(id int, userId int) error {
	query := `
		DELETE FROM categories WHERE id = $1 AND user_id = $2
	`
	result, err := r.db.Exec(query, id, userId)
	if err != nil {
		return fmt.Errorf("failed to delete category by id")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("category not found")
	}
	return nil
}
