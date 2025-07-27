package models

type Categories struct {
	ID           int64  `json:"id"`
	UserId       int64  `json:"user_id"`
	CategoryName string `json:"category_name"`
	CreateAt     string `json:"create_at"`
	CreateBy     int64  `json:"create_by"`
	ModifiedAt   string `json:"modified_at"`
	ModifiedBy   int64  `json:"modified_by"`
}
