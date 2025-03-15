package models

import "time"

type LoanInvestment struct {
	ID              int64     `json:"id"`
	LoanID          int64     `json:"loan_id" binding:"required" form:"loan_id"`
	InvestorID      int64     `json:"investor_id" binding:"required" form:"investor_id"`
	AgreementLetter string    `json:"agreement_letter"`
	InvestedAmount  float64   `json:"invested_amount" binding:"required" form:"invested_amount"`
	CreatedAt       time.Time `json:"created_at"`
}

type InvestorEmail struct {
	InvestorID      int64  `json:"investor_id"`
	InvestorName    string `json:"investor_name"`
	InvestorEmail   string `json:"investor_email"`
	AgreementLetter string `json:"agreement_letter"`
}
