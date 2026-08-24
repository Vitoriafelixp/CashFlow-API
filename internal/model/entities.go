package model

import (
	"time"
)

type Lancamento struct {
	ID                      int64      `db:"id" json:"id"`
	Descricao               string     `db:"description" json:"descricao"`
	Valor                   float64    `db:"amount" json:"valor"`
	Tipo                    string     `db:"type" json:"tipo"`
	DataVencimento          *time.Time `db:"due_date" json:"data_vencimento,omitempty"`
	DataPrevistaRecebimento *time.Time `db:"expected_receipt_date" json:"data_prevista_recebimento,omitempty"`
	Categoria               string     `db:"category" json:"categoria"`
	Status                  string     `db:"status" json:"status"`
	DataPagamento           *time.Time `db:"paid_at" json:"data_pagamento,omitempty"`
	Confirmado              *bool      `db:"confirmed" json:"confirmado,omitempty"`
	CreatedAt               time.Time  `db:"created_at" json:"created_at"`
}
type Configuracoes struct {
	ID            int     `db:"id" json:"id"`
	OrcamentoReposicao float64 `db:"restock_budget" json:"orcamento_reposicao"`
}
