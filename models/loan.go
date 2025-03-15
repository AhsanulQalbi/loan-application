package models

type Loan struct {
	ID               int64   `json:"id"`
	BorrowerID       int64   `json:"borrower_id" binding:"required" form:"borrower_id"`
	BorrowerUserName string  `json:"borrower_user_name"`
	Principal        float64 `json:"principal" binding:"required" form:"principal_amount"`
	TotalInvested    float64 `json:"total_invested"`
	Rate             float64 `json:"rate" binding:"required" form:"rate"`
	ROI              float64 `json:"roi" binding:"required" form:"roi"`
	LoanState        string  `json:"loan_state"`
	CreatedAt        string  `json:"created_at"`
}

type UpdateLoanRequest struct {
	PrincipalAmount *float64 `json:"principal_amount" form:"principal_amount"`
	Rate            *float64 `json:"rate" form:"rate"`
	ROI             *float64 `json:"roi" form:"roi"`
}
