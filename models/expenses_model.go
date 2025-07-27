package models

import (
	"exp_tracker/utils"
	"time"
)

type Expenses struct {
	ID            int64          `json:"id"`
	UserId        int64          `json:"user_id"`
	CategoryId    int64          `json:"category_id"`
	PaymentMethod string         `json:"payment_method"`
	Title         string         `json:"title"`
	Amount        int            `json:"amount"`
	Description   string         `json:"description"`
	ExpenseDate   utils.DateOnly `json:"expense_date"`
	CreateAt      time.Time      `json:"create_at"`
	CreateBy      int64          `json:"create_by"`
	ModifiedAt    time.Time      `json:"modified_at"`
	ModifiedBy    int64          `json:"modified_by"`
}

type DateRangeExpenses struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}
