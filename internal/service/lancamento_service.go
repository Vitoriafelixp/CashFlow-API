package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"cashflow-manager/internal/model"
	"cashflow-manager/internal/repository"
)

var ErrInvalidLancamento = errors.New("lancamento invalido")

type LancamentoService struct {
	repository *repository.LancamentoRepository
}

func NewLancamentoService(repository *repository.LancamentoRepository) *LancamentoService {
	return &LancamentoService{repository: repository}
}

func (s *LancamentoService) Criar(ctx context.Context, lancamento model.Entry) (model.Entry, error) {
	if s == nil || s.repository == nil {
		return model.Entry{}, fmt.Errorf("lancamento repository nao inicializado")
	}

	lancamento.Description = strings.TrimSpace(lancamento.Description)
	lancamento.Category = strings.TrimSpace(lancamento.Category)
	lancamento.Type = strings.TrimSpace(lancamento.Type)

	if lancamento.Description == "" {
		return model.Entry{}, fmt.Errorf("%w: descricao e obrigatoria", ErrInvalidLancamento)
	}
	if lancamento.Amount <= 0 {
		return model.Entry{}, fmt.Errorf("%w: valor deve ser maior que zero", ErrInvalidLancamento)
	}

	switch lancamento.Type {
	case "receivable":
		if lancamento.ExpectedReceiptDate == nil || lancamento.DueDate != nil {
			return model.Entry{}, fmt.Errorf("%w: recebiveis exigem expect_receipt_date e nao aceitam due_date", ErrInvalidLancamento)
		}
	case "payable":
		if lancamento.DueDate == nil || lancamento.ExpectedReceiptDate != nil {
			return model.Entry{}, fmt.Errorf("%w: despesas exigem due_date e nao aceitam expect_receipt_date", ErrInvalidLancamento)
		}
	default:
		return model.Entry{}, fmt.Errorf("%w: type deve ser receivable ou payable", ErrInvalidLancamento)
	}

	lancamento.Status = "pending"
	lancamento.PaidAt = nil

	return s.repository.Criar(ctx, lancamento)
}

func (s *LancamentoService) Listar(ctx context.Context) ([]model.Entry, error) {
	if s == nil || s.repository == nil {
		return nil, fmt.Errorf("lancamento repository nao inicializado")
	}

	return s.repository.Listar(ctx, repository.FiltroLancamento{})
}
