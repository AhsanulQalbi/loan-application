package repositories

import (
	"database/sql"
	"fmt"
	"loan_application/models"
)

type LoanRepository struct {
	DB        *sql.DB
	stateRepo *LoanStateRepository
}

func NewLoanRepository(db *sql.DB, stateRepo *LoanStateRepository) *LoanRepository {
	return &LoanRepository{DB: db, stateRepo: stateRepo}
}

func (r *LoanRepository) GetLoans() ([]models.Loan, error) {
	rows, err := r.DB.Query(`
		SELECT 
			loans.id, 
			loans.borrower_id, 
			borrowers.user_name, 
			loans.principal_amount, 
			COALESCE(SUM(loan_investments.invested_amount), 0) AS total_invested, 
			loans.rate, 
			loans.roi, 
			loans.loan_state, 
			loans.created_at 
		FROM loans
		JOIN borrowers ON borrowers.id = loans.borrower_id
		LEFT JOIN loan_investments ON loans.id = loan_investments.loan_id
		GROUP BY 
			loans.id, 
			borrowers.user_name,
			loans.borrower_id, 
			loans.principal_amount, 
			loans.rate, 
			loans.roi, 
			loans.loan_state, 
			loans.created_at
		ORDER BY loans.id ASC;
	`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var loans []models.Loan
	for rows.Next() {
		var loan models.Loan
		if err := rows.Scan(&loan.ID, &loan.BorrowerID, &loan.BorrowerUserName, &loan.Principal, &loan.TotalInvested, &loan.Rate, &loan.ROI, &loan.LoanState, &loan.CreatedAt); err != nil {
			return nil, err
		}

		loans = append(loans, loan)
	}

	return loans, nil
}

func (r *LoanRepository) GetLoanByID(id int64) (*models.Loan, error) {
	var loan models.Loan
	query := "SELECT id, borrower_id, principal_amount, rate, roi, loan_state, created_at FROM loans WHERE id = $1"
	err := r.DB.QueryRow(query, id).Scan(&loan.ID, &loan.BorrowerID, &loan.Principal, &loan.Rate, &loan.ROI, &loan.LoanState, &loan.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &loan, nil
}

func (r *LoanRepository) CreateLoan(loan *models.Loan) (int64, error) {
	query := "INSERT INTO loans (borrower_id, principal_amount, rate, roi, loan_state) VALUES ($1, $2, $3, $4, 'proposed') RETURNING id"
	err := r.DB.QueryRow(query, loan.BorrowerID, loan.Principal, loan.Rate, loan.ROI).Scan(&loan.ID)
	if err != nil {
		return 0, err
	}

	if err := r.stateRepo.CreateLoanState(loan.ID, "proposed"); err != nil {
		return 0, fmt.Errorf("failed to record loan state: %w", err)
	}

	return loan.ID, nil
}

func (r *LoanRepository) UpdateLoan(loanID int64, principalAmount, rate, roi *float64) error {
	query := `
		UPDATE loans 
		SET 
			principal_amount = COALESCE($1, principal_amount),
			rate = COALESCE($2, rate),
			roi = COALESCE($3, roi)
		WHERE id = $4
	`
	_, err := r.DB.Exec(query, principalAmount, rate, roi, loanID)
	if err != nil {
		return fmt.Errorf("error updating loan: %w", err)
	}
	return nil
}
