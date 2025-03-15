package repositories

import (
	"database/sql"
	"loan_application/models"
)

type BorrowerRepository struct {
	DB *sql.DB
}

func NewBorrowerRepository(db *sql.DB) *BorrowerRepository {
	return &BorrowerRepository{DB: db}
}

func (r *BorrowerRepository) GetBorrowers() ([]models.Borrower, error) {
	rows, err := r.DB.Query("SELECT id, user_name, email, created_at FROM borrowers")
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	var borrowers []models.Borrower
	for rows.Next() {
		var borrower models.Borrower
		if err := rows.Scan(&borrower.ID, &borrower.Username, &borrower.Email, &borrower.CreatedAt); err != nil {
			return nil, err
		}

		borrowers = append(borrowers, borrower)
	}

	return borrowers, nil
}

func (r *BorrowerRepository) GetBorrowerByID(id int64) (*models.Borrower, error) {
	var borrower models.Borrower
	query := "SELECT id, user_name, created_at FROM borrowers WHERE id = $1"
	err := r.DB.QueryRow(query, id).Scan(&borrower.ID, &borrower.Username, &borrower.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &borrower, nil
}

func (r *BorrowerRepository) CreateBorrower(borrower *models.Borrower) (int64, error) {
	query := "INSERT INTO borrowers (user_name, email, created_at) VALUES ($1,$2,NOW()) RETURNING id"
	err := r.DB.QueryRow(query, borrower.Username, borrower.Email).Scan(&borrower.ID)
	if err != nil {
		return 0, err
	}

	return borrower.ID, nil
}
