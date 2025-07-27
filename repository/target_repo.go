package repository

import (
	"database/sql"
	"exp_tracker/models"
	"fmt"
)

type TargetRepository interface {
	Create(target *models.Target) error
	GetTargetExpense(userId string) ([]models.Target, error)
	UpdateTarget(id int, target models.Target) (*models.Target, error)
	GetTotalAmountByUser(userId int) (int64, error)
	DeleteTargetByID(id int64) error
}

type TargetRepo struct {
	db *sql.DB
}

func CTargetRepository(db *sql.DB) TargetRepository {
	return &TargetRepo{db}
}

func (r *TargetRepo) Create(b *models.Target) error {
	query := `
		INSERT INTO target (
		id, 
		user_id, 
		file, 
		title, 
		payment_method, 
		description, 
		amount, 
		total_amount,
		start_date, 
		end_date, 
		create_at, 
		create_by, 
		modified_at, 
		modified_by
		) VALUES (
		  $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
		)
	`
	_, err := r.db.Exec(
		query,
		b.ID,
		b.UserId,
		b.File,
		b.Title,
		b.PaymentMethod,
		b.Description,
		b.Amount,
		b.TotalAmount,
		b.StartDate,
		b.EndDate,
		b.CreateAt,
		b.CreateBy,
		b.ModifiedAt,
		b.ModifiedBy)

	if err != nil {
		return fmt.Errorf("failed to insert target expense %w", err)
	}
	return nil
}

func (r *TargetRepo) GetTotalAmountByUser(userId int) (int64, error) {
	query := `SELECT COALESCE(SUM(amount), 0) FROM target WHERE user_id = $1`
	var total int64
	err := r.db.QueryRow(query, userId).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to get total amount: %w", err)
	}
	return total, nil
}

func (r *TargetRepo) GetTargetExpense(userId string) ([]models.Target, error) {
	query := `
		SELECT * FROM target WHERE user_id = $1
	`
	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get all expenses: %w", err)
	}
	defer rows.Close()

	var targets []models.Target

	for rows.Next() {
		var t models.Target
		err := rows.Scan(
			&t.ID,
			&t.UserId,
			&t.File,
			&t.Title,
			&t.PaymentMethod,
			&t.Description,
			&t.Amount,
			&t.TotalAmount,
			&t.StartDate,
			&t.EndDate,
			&t.CreateAt,
			&t.CreateBy,
			&t.ModifiedAt,
			&t.ModifiedBy)
		if err != nil {
			return nil, fmt.Errorf("failed to scan expense %w", err)
		}
		targets = append(targets, t)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return targets, nil
}

func (r *TargetRepo) UpdateTarget(id int, target models.Target) (*models.Target, error) {
	query := `
		UPDATE target
		SET title = $2,
			file = $3,
			payment_method = $4,
			description = $5,
			amount = $6,
			total_amount = $7,
			start_date = $8,
			end_date = $9,
			modified_at = $10,
			modified_by = $11
		WHERE id = $1
		RETURNING *
	`
	rows, err := r.db.Exec(
		query,
		id,
		target.Title,
		target.File,
		target.PaymentMethod,
		target.Description,
		target.Amount,
		target.TotalAmount,
		target.StartDate,
		target.EndDate,
		target.ModifiedAt,
		target.ModifiedBy,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update expense %w", err)
	}

	rowsAffected, err := rows.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("rows affected failed %w", err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("target expense not found")
	}

	return &target, nil
}

func (r *TargetRepo) DeleteTargetByID(id int64) error {
	query := `
		DELETE FROM target WHERE id = $1
	`

	rows, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete target %w", err)
	}

	rowsAffected, err := rows.RowsAffected()
	if err != nil {

		if rowsAffected == 0 {
			return fmt.Errorf("target not found %w", err)
		}
		return fmt.Errorf("error %w", err)
	}
	return nil
}
