package models

import "time"

type LoanState struct {
	ID        int       `json:"id"`
	LoanID    int       `json:"loan_id"`
	State     string    `json:"state"`
	ChangedAt time.Time `json:"changed_at"`
}
