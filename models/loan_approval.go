package models

import "time"

type LoanApproval struct {
	ID                  int64     `json:"id"`
	LoanID              int64     `json:"loan_id" binding:"required" form:"loan_id"`
	EmployeeValidatorID int64     `json:"employee_validator_id" binding:"required" form:"employee_validator_id"`
	VisitProof          string    `json:"visit_proof"`
	ApprovalDate        time.Time `json:"approval_date"`
	CreatedAt           time.Time `json:"created_at"`
}
