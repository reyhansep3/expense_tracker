package models

import "time"

type Target struct {
	ID            int64     `json:"id"`
	UserId        int64     `json:"user_id"`
	File          string    `json:"file"`
	Title         string    `json:"title"`
	PaymentMethod string    `json:"payment_method"`
	Description   string    `json:"description"`
	Amount        int64     `json:"amount"`
	TotalAmount   int64     `json:"total_amount"`
	StartDate     time.Time `json:"start_date"`
	EndDate       time.Time `json:"end_date"`
	CreateAt      time.Time `json:"create_at"`
	CreateBy      int64     `json:"create_by"`
	ModifiedAt    time.Time `json:"modified_at"`
	ModifiedBy    int64     `json:"modified_by"`
}
