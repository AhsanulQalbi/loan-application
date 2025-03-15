package models

import "time"

type LoanDisbursement struct {
	ID                int64     `json:"id"`
	LoanID            int64     `json:"loan_id" binding:"required" form:"loan_id"`
	AgreementLetter   string    `json:"agreement_letter"`
	EmployeeOfficerID int64     `json:"employee_officer_id" binding:"required" form:"employee_officer_id"`
	DisbursementDate  time.Time `json:"disbursement_date"`
	CreatedAt         time.Time `json:"created_at"`
}
