package repositories

import (
	"database/sql"
	"fmt"
	"loan_application/models"
	"log"
	"net/smtp"
	"os"
	"strings"
)

type LoanInvestmentRepository struct {
	DB            *sql.DB
	LoanStateRepo *LoanStateRepository
	InvestorRepo  *InvestorRepository
}

func NewLoanInvestmentRepository(db *sql.DB, loanStateRepo *LoanStateRepository, investorRepo *InvestorRepository) *LoanInvestmentRepository {
	return &LoanInvestmentRepository{
		DB:            db,
		LoanStateRepo: loanStateRepo,
		InvestorRepo:  investorRepo,
	}
}

func (repo *LoanInvestmentRepository) InvestInLoan(investment models.LoanInvestment) (float64, error) {
	tx, err := repo.DB.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to start transaction: %w", err)
	}

	var principalAmount float64
	query := `SELECT principal_amount FROM loans WHERE id = $1`
	err = tx.QueryRow(query, investment.LoanID).Scan(&principalAmount)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("error fetching loan principal amount: %w", err)
	}

	var totalInvested float64
	query = `SELECT COALESCE(SUM(invested_amount), 0) FROM loan_investments WHERE loan_id = $1`
	err = tx.QueryRow(query, investment.LoanID).Scan(&totalInvested)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("error calculating total invested: %w", err)
	}

	if totalInvested+investment.InvestedAmount > principalAmount {
		tx.Rollback()
		return 0, fmt.Errorf("invested amount will exceeds loan principal amount")
	}

	query = `INSERT INTO loan_investments (loan_id, investor_id, agreement_letter, invested_amount, created_at)
	          VALUES ($1, $2, $3, $4, NOW())`
	_, err = tx.Exec(query, investment.LoanID, investment.InvestorID, investment.AgreementLetter, investment.InvestedAmount)
	if err != nil {
		tx.Rollback()
		if strings.Contains(err.Error(), "pq: duplicate key value violates") {
			return 0, fmt.Errorf("error inserting investment: this investor already invested in this loan")
		}
		return 0, fmt.Errorf("error inserting investment: %w", err)
	}

	totalInvested += investment.InvestedAmount
	if totalInvested >= principalAmount {
		query = `UPDATE loans SET loan_state = 'invested' WHERE id = $1`
		_, err = tx.Exec(query, investment.LoanID)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("error updating loan state to invested: %w", err)
		}

		if err := repo.LoanStateRepo.CreateLoanState(investment.LoanID, "invested"); err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to record loan state: %w", err)
		}

		investors, err := repo.InvestorRepo.GetInvestorsByLoanID(investment.LoanID)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("failed to get investors by loanID: %w", err)
		}

		err = tx.Commit()
		if err != nil {
			return 0, fmt.Errorf("failed to commit transaction: %w", err)
		}

		go func() {
			for _, investor := range investors {
				SendInvestmentEmail(investor)
			}
		}()

	} else {
		err = tx.Commit()
		if err != nil {
			return 0, fmt.Errorf("failed to commit transaction: %w", err)
		}
	}

	return totalInvested, nil
}

func SendInvestmentEmail(investor models.InvestorEmail) {
	from := os.Getenv("EMAIL_FROM")
	password := os.Getenv("EMAIL_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	to := []string{investor.InvestorEmail}
	subject := "Your Investment Agreement"
	body := fmt.Sprintf("Dear %s,\n\nThank you for your investment. You can download your agreement letter here: %s\n\nBest regards,\nLoan Platform Team",
		investor.InvestorName, investor.AgreementLetter)

	message := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", from, investor.InvestorEmail, subject, body)

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, []byte(message))
	if err != nil {
		log.Println("Failed to send email to", investor.InvestorEmail, ":", err)
	} else {
		log.Println("Email sent successfully to", investor.InvestorEmail)
	}
}
