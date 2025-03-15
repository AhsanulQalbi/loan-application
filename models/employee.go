package models

type Employee struct {
	ID        int64  `json:"id"`
	Email     string `json:"email" binding:"required" form:"email"`
	Username  string `json:"user_name" binding:"required" form:"user_name"`
	CreatedAt string `json:"created_at"`
}
