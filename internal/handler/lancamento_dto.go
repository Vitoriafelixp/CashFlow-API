package handler

import (
	"time"

	"cashflow-manager/internal/model"
)

type CriarLancamentoRequest struct {
	Description         string     `json:"description"`
	Amount              float64    `json:"amount"`
	Type                string     `json:"type"`
	DueDate             *time.Time `json:"due_date,omitempty"`
	ExpectedReceiptDate *time.Time `json:"expect_receipt_date,omitempty"`
	Category            string     `json:"category"`
	Confirmed           *bool      `json:"confirmed,omitempty"`
}

type LancamentoResponse struct {
	ID                  int        `json:"id"`
	Description         string     `json:"description"`
	Amount              float64    `json:"amount"`
	Type                string     `json:"type"`
	DueDate             *time.Time `json:"due_date,omitempty"`
	ExpectedReceiptDate *time.Time `json:"expect_receipt_date,omitempty"`
	Category            string     `json:"category"`
	Confirmed           *bool      `json:"confirmed,omitempty"`
	Status              string     `json:"status"`
	PaidAt              *time.Time `json:"paid_at,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
}

type ListaLancamentosResponse []LancamentoResponse

func (request CriarLancamentoRequest) toModel() model.Entry {
	return model.Entry{
		Description:         request.Description,
		Amount:              request.Amount,
		Type:                request.Type,
		DueDate:             request.DueDate,
		ExpectedReceiptDate: request.ExpectedReceiptDate,
		Category:            request.Category,
		Confirmed:           request.Confirmed,
	}
}

func toLancamentoResponse(lancamento model.Entry) LancamentoResponse {
	return LancamentoResponse{
		ID:                  lancamento.ID,
		Description:         lancamento.Description,
		Amount:              lancamento.Amount,
		Type:                lancamento.Type,
		DueDate:             lancamento.DueDate,
		ExpectedReceiptDate: lancamento.ExpectedReceiptDate,
		Category:            lancamento.Category,
		Confirmed:           lancamento.Confirmed,
		Status:              lancamento.Status,
		PaidAt:              lancamento.PaidAt,
		CreatedAt:           lancamento.CreatedAt,
	}
}

func toListaLancamentosResponse(lancamentos []model.Entry) ListaLancamentosResponse {
	resposta := make(ListaLancamentosResponse, 0, len(lancamentos))
	for _, lancamento := range lancamentos {
		resposta = append(resposta, toLancamentoResponse(lancamento))
	}
	return resposta
}
