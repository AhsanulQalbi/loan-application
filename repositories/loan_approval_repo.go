package repositories

import (
	"database/sql"
	"fmt"
	"loan_application/models"
)

type LoanApprovalRepository struct {
	DB            *sql.DB
	loanStateRepo *LoanStateRepository
}

func NewLoanApprovalRepository(db *sql.DB, loanStateRepo *LoanStateRepository) *LoanApprovalRepository {
	return &LoanApprovalRepository{
		DB:            db,
		loanStateRepo: loanStateRepo,
	}
}

func (r *LoanApprovalRepository) GetLoanApprovalByLoanID(id int64) (*models.LoanApproval, error) {
	var loanApproval models.LoanApproval
	query := "SELECT id, loan_id, employee_validator_id, visit_proof, approval_date FROM loans WHERE loan_id = $1"
	err := r.DB.QueryRow(query, id).Scan(&loanApproval.ID, &loanApproval.EmployeeValidatorID, &loanApproval.VisitProof, &loanApproval.ApprovalDate)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &loanApproval, nil
}

func (repo *LoanApprovalRepository) ApproveLoan(approval models.LoanApproval) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	query := `INSERT INTO loan_approvals (loan_id, employee_validator_id, visit_proof, approval_date, created_at)
	          VALUES ($1, $2, $3, NOW(), NOW())`
	_, err = tx.Exec(query, approval.LoanID, approval.EmployeeValidatorID, approval.VisitProof)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error inserting into loan_approvals: %w", err)
	}

	query = `UPDATE loans SET loan_state = 'approved' WHERE id = $1`
	_, err = tx.Exec(query, approval.LoanID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error updating loan state: %w", err)
	}

	if err := repo.loanStateRepo.CreateLoanState(approval.LoanID, "approved"); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to record loan state: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
