package model

import (
	"time"
)

type Entry struct {
	ID                  int        `db:"id" json:"id"`
	Description         string     `db:"description" json:"description"`
	Amount              float64    `db:"amount" json:"amount"`
	Type                string     `db:"type" json:"type"`
	DueDate             *time.Time `db:"due_date" json:"due_date,omitempty"`
	ExpectedReceiptDate *time.Time `db:"expect_receipt_date" json:"expect_receipt_date,omitempty"`
	Category            string     `db:"category" json:"category"`
	Confirmed           *bool      `db:"confirmed" json:"confirmed,omitempty"`
	Status              string     `db:"status" json:"status"`
	PaidAt              *time.Time `db:"paid_at" json:"paid_at,omitempty"`
	CreatedAt           time.Time  `db:"created_at" json:"created_at"`
}
type Settings struct {
	ID            int     `db:"id" json:"id"`
	RestockBudget float64 `db:"restock_budget" json:"restock_budget"`
}
