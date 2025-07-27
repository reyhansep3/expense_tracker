package repository

import (
	"database/sql"
	"exp_tracker/models"
	"fmt"
	"time"
)

type ExpenseRepository interface {
	Create(expense *models.Expenses) error
	GetUserExpense(userId int) ([]models.Expenses, error)
	GetExpenseByDate(userId int, startDate, endDate time.Time) ([]models.Expenses, error)
	UpdateUserExpense(id int, expense models.Expenses, expenseDate time.Time) (*models.Expenses, error)
	GetTotalExpenseByUserWithDateRange(userId int, startDate, endDate *time.Time) (int, error)
	DeleteExpenseByID(id int64) error
}

type ExpenseRepo struct {
	db *sql.DB
}

func CExpenseRepository(db *sql.DB) ExpenseRepository {
	return &ExpenseRepo{db}
}

func (r *ExpenseRepo) Create(b *models.Expenses) error {
	query := `
	INSERT INTO expenses (
		id, user_id, category_id, payment_method, title, amount, description, expense_date, create_at, create_by, modified_at, modified_by
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12
	)
	`

	_, err := r.db.Exec(query,
		b.ID,
		b.UserId,
		b.CategoryId,
		b.PaymentMethod,
		b.Title,
		b.Amount,
		b.Description,
		b.ExpenseDate.ToTime(),
		b.CreateAt,
		b.CreateBy,
		b.ModifiedAt,
		b.ModifiedBy,
	)
	if err != nil {
		return fmt.Errorf("failed to insert expenses: %w", err)
	}
	return nil
}

func (r *ExpenseRepo) GetUserExpense(userId int) ([]models.Expenses, error) {
	query := `
		SELECT * FROM expenses WHERE user_id = $1
	`

	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get all expenses: %w", err)
	}

	var expenses []models.Expenses

	for rows.Next() {
		var e models.Expenses
		err := rows.Scan(&e.ID, &e.UserId, &e.CategoryId, &e.PaymentMethod, &e.Title, &e.Amount, &e.Description, &e.ExpenseDate, &e.CreateAt, &e.CreateBy, &e.ModifiedAt, &e.ModifiedBy)
		if err != nil {
			return nil, fmt.Errorf("failed to scan expense %w", err)
		}
		expenses = append(expenses, e)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return expenses, nil

}

func (r *ExpenseRepo) GetExpenseByDate(userId int, startDate, endDate time.Time) ([]models.Expenses, error) {
	query := `
		SELECT * FROM expenses 
		WHERE user_id = $1 
		AND expense_date 
		BETWEEN $2 AND $3
	`
	var expenses []models.Expenses
	rows, err := r.db.Query(query, userId, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get expense by date %w", err)
	}
	for rows.Next() {
		var e models.Expenses

		err := rows.Scan(&e.ID, &e.UserId, &e.CategoryId, &e.PaymentMethod, &e.Title, &e.Amount, &e.Description, &e.ExpenseDate, &e.CreateAt, &e.CreateBy, &e.ModifiedAt, &e.ModifiedBy)
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, fmt.Errorf("expenses not found %w", err)
			}
			return nil, fmt.Errorf("failed to get expense by id %w", err)
		}
		expenses = append(expenses, e)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return expenses, nil
}

func (r *ExpenseRepo) UpdateUserExpense(id int, expense models.Expenses, expenseDate time.Time) (*models.Expenses, error) {
	query := `
        UPDATE expenses
        SET category_id = $2,
            payment_method = $3,
            title = $4,
            amount = $5,
            description = $6,
            expense_date = $7,
            modified_at = $8,
            modified_by = $9
        WHERE id = $1
    `

	res, err := r.db.Exec(
		query, id,
		expense.CategoryId,
		expense.PaymentMethod,
		expense.Title,
		expense.Amount,
		expense.Description,
		expenseDate,
		expense.ModifiedAt,
		expense.ModifiedBy,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update expense %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get affected rows %w", err)
	}
	if rowsAffected == 0 {
		return nil, fmt.Errorf("expense not found")
	}

	return &expense, nil
}

func (r *ExpenseRepo) DeleteExpenseByID(id int64) error {
	query := `
		DELETE FROM expenses WHERE id = $1
	`

	rows, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete expenses %w", err)
	}

	rowsAffected, err := rows.RowsAffected()
	if err != nil {

		if rowsAffected == 0 {
			return fmt.Errorf("expense not found %w", err)
		}
		return fmt.Errorf("error %w", err)
	}
	return nil
}

func (r *ExpenseRepo) GetTotalExpenseByUserWithDateRange(userId int, startDate, endDate *time.Time) (int, error) {
	query := `
		SELECT COALESCE(SUM(amount), 0)
		FROM expenses
		WHERE user_id = $1
	`

	args := []interface{}{userId}
	paramIndex := 2

	if startDate != nil {
		query += fmt.Sprintf(" AND expense_date >= $%d", paramIndex)
		args = append(args, *startDate)
		paramIndex++
	}
	if endDate != nil {
		query += fmt.Sprintf(" AND expense_date <= $%d", paramIndex)
		args = append(args, *endDate)
		paramIndex++
	}

	var total int
	err := r.db.QueryRow(query, args...).Scan(&total)
	if err != nil {
		return 0, fmt.Errorf("failed to get total expenses with date range: %w", err)
	}
	return total, nil
}
