package models

type Investor struct {
	ID        int64  `json:"id"`
	Username  string `json:"user_name" form:"user_name" binding:"required" `
	Email     string `json:"email" form:"email" binding:"required"`
	CreatedAt string `json:"created_at"`
}
