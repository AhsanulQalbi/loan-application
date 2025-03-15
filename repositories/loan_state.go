package repositories

import (
	"database/sql"
	"fmt"
	"loan_application/models"
)

type LoanStateRepository struct {
	DB *sql.DB
}

func NewLoanStateRepository(db *sql.DB) *LoanStateRepository {
	return &LoanStateRepository{
		DB: db,
	}
}

func (repo *LoanStateRepository) CreateLoanState(loanID int64, state string) error {
	query := `INSERT INTO loan_states (loan_id, state, changed_at) VALUES ($1, $2, NOW())`
	_, err := repo.DB.Exec(query, loanID, state)
	if err != nil {
		return fmt.Errorf("error inserting loan state: %w", err)
	}

	return nil
}

func (r *LoanStateRepository) GetLoanStates() ([]models.LoanState, error) {
	rows, err := r.DB.Query("SELECT id, loan_id, state, changed_at FROM loan_states order by 1 desc")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var loanStates []models.LoanState
	for rows.Next() {
		var loanState models.LoanState
		if err := rows.Scan(&loanState.ID, &loanState.LoanID, &loanState.State, &loanState.ChangedAt); err != nil {
			return nil, err
		}

		loanStates = append(loanStates, loanState)
	}

	return loanStates, nil
}
