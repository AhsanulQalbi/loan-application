package repositories

import (
	"database/sql"
	"loan_application/models"
)

type InvestorRepository struct {
	DB *sql.DB
}

func NewInvestorRepository(db *sql.DB) *InvestorRepository {
	return &InvestorRepository{DB: db}
}

func (i *InvestorRepository) GetInvestors() ([]models.Investor, error) {
	rows, err := i.DB.Query("SELECT id, user_name, email, created_at FROM investors")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var investors []models.Investor
	for rows.Next() {
		var investor models.Investor
		if err := rows.Scan(&investor.ID, &investor.Username, &investor.Email, &investor.CreatedAt); err != nil {
			return nil, err
		}

		investors = append(investors, investor)
	}

	return investors, nil
}

func (r *InvestorRepository) GetInvestorByID(id int64) (*models.Investor, error) {
	var investor models.Investor
	query := "SELECT id, user_name, created_at FROM investors WHERE id = $1"
	err := r.DB.QueryRow(query, id).Scan(&investor.ID, &investor.Username, &investor.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &investor, nil
}

func (r *InvestorRepository) CreateInvestors(investor *models.Investor) (int64, error) {
	query := "INSERT INTO investors (user_name, email, created_at) VALUES ($1, $2, NOW()) RETURNING id"
	err := r.DB.QueryRow(query, investor.Username, investor.Email).Scan(&investor.ID)
	if err != nil {
		return 0, err
	}

	return investor.ID, nil
}

func (repo *InvestorRepository) GetInvestorsByLoanID(loanID int64) ([]models.InvestorEmail, error) {
	query := `
	SELECT 
	    investors.id, 
	    investors.user_name, 
	    investors.email, 
	    loan_investments.agreement_letter
	FROM loan_investments
	JOIN investors ON investors.id = loan_investments.investor_id
	WHERE loan_investments.loan_id = $1;
	`

	rows, err := repo.DB.Query(query, loanID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var investors []models.InvestorEmail
	for rows.Next() {
		var investor models.InvestorEmail
		if err := rows.Scan(&investor.InvestorID, &investor.InvestorName, &investor.InvestorEmail, &investor.AgreementLetter); err != nil {
			return nil, err
		}
		investors = append(investors, investor)
	}

	return investors, nil
}
