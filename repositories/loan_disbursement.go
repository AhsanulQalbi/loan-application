package repositories

import (
	"database/sql"
	"fmt"
	"loan_application/models"
)

type LoanDisbursementRepository struct {
	DB            *sql.DB
	LoanStateRepo *LoanStateRepository
}

func NewLoanDisbirsementRepository(db *sql.DB, loanStateRepo *LoanStateRepository) *LoanDisbursementRepository {
	return &LoanDisbursementRepository{
		DB:            db,
		LoanStateRepo: loanStateRepo,
	}
}

func (repo *LoanDisbursementRepository) DisburseLoan(disbursement models.LoanDisbursement) error {
	tx, err := repo.DB.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	var currentState string
	query := `SELECT loan_state FROM loans WHERE id = $1`
	err = tx.QueryRow(query, disbursement.LoanID).Scan(&currentState)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error get loan loan state: %w", err)
	}

	if currentState != "invested" {
		tx.Rollback()
		return fmt.Errorf("loan is not fully invested yet")
	}

	query = `INSERT INTO loan_disbursements (loan_id, agreement_letter, employee_officer_id, disbursement_date, created_at)
	         VALUES ($1, $2, $3, NOW(), NOW())`
	_, err = tx.Exec(query, disbursement.LoanID, disbursement.AgreementLetter, disbursement.EmployeeOfficerID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error inserting loan disbursement: %w", err)
	}

	query = `UPDATE loans SET loan_state = 'disbursed' WHERE id = $1`
	_, err = tx.Exec(query, disbursement.LoanID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error updating loan state to disbursed: %w", err)
	}

	if err := repo.LoanStateRepo.CreateLoanState(disbursement.LoanID, "disbursed"); err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to record loan state: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("error, failed to commit transaction: %w", err)
	}

	return nil
}
