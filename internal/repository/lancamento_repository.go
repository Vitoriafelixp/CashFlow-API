package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"cashflow-manager/internal/database"
	"cashflow-manager/internal/model"
)

type FiltroLancamento struct {
	Status                        string
	Tipo                          string
	DataVencimentoInicio          *time.Time
	DataVencimentoFim             *time.Time
	DataPrevistaRecebimentoInicio *time.Time
	DataPrevistaRecebimentoFim    *time.Time
	Confirmado                    *bool
}

type LancamentoRepository struct {
	db *database.DB
}

func NewLancamentoRepository(db *database.DB) *LancamentoRepository {
	return &LancamentoRepository{db: db}
}

func (r *LancamentoRepository) Criar(ctx context.Context, lancamento model.Entry) (model.Entry, error) {
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO entries (
			description, amount, type, due_date, expect_receipt_date,
			category, confirmed, status
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`,
		lancamento.Description,
		lancamento.Amount,
		lancamento.Type,
		lancamento.DueDate,
		lancamento.ExpectedReceiptDate,
		lancamento.Category,
		lancamento.Confirmed,
		lancamento.Status,
	)
	if err != nil {
		return model.Entry{}, fmt.Errorf("criar lancamento: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return model.Entry{}, fmt.Errorf("obter id do lancamento criado: %w", err)
	}

	return r.ObterPorID(ctx, int(id))
}

func (r *LancamentoRepository) ObterPorID(ctx context.Context, id int) (model.Entry, error) {
	var lancamento model.Entry
	err := r.db.GetContext(ctx, &lancamento, `
		SELECT id, description, amount, type, due_date, expect_receipt_date,
			COALESCE(category, '') AS category, confirmed, status, paid_at, created_at
		FROM entries
		WHERE id = ?
	`, id)
	if err != nil {
		return model.Entry{}, fmt.Errorf("obter lancamento criado: %w", err)
	}

	return lancamento, nil
}

func (r *LancamentoRepository) Listar(ctx context.Context, filtro FiltroLancamento) ([]model.Entry, error) {
	query := `
		SELECT id, description, amount, type, due_date, expect_receipt_date,
			COALESCE(category, '') AS category, confirmed, status, paid_at, created_at
		FROM entries
		WHERE 1 = 1
	`
	args := make([]any, 0)

	if filtro.Status != "" {
		query += " AND status = ?"
		args = append(args, filtro.Status)
	}
	if filtro.Tipo != "" {
		query += " AND type = ?"
		args = append(args, filtro.Tipo)
	}
	if filtro.DataVencimentoInicio != nil {
		query += " AND due_date >= ?"
		args = append(args, filtro.DataVencimentoInicio.Format("2006-01-02"))
	}
	if filtro.DataVencimentoFim != nil {
		query += " AND due_date <= ?"
		args = append(args, filtro.DataVencimentoFim.Format("2006-01-02"))
	}
	if filtro.DataPrevistaRecebimentoInicio != nil {
		query += " AND expect_receipt_date >= ?"
		args = append(args, filtro.DataPrevistaRecebimentoInicio.Format("2006-01-02"))
	}
	if filtro.DataPrevistaRecebimentoFim != nil {
		query += " AND expect_receipt_date <= ?"
		args = append(args, filtro.DataPrevistaRecebimentoFim.Format("2006-01-02"))
	}
	if filtro.Confirmado != nil {
		query += " AND confirmed = ?"
		args = append(args, *filtro.Confirmado)
	}

	query += " ORDER BY due_date ASC, id ASC"

	lancamentos := make([]model.Entry, 0)
	if err := r.db.SelectContext(ctx, &lancamentos, query, args...); err != nil {
		return nil, fmt.Errorf("listar lancamentos: %w", err)
	}

	return lancamentos, nil
}

func (r *LancamentoRepository) ListarVencidos(ctx context.Context) ([]model.Entry, error) {
	inicio := time.Now().AddDate(0, 0, -30)
	fim := time.Now().AddDate(0, 0, 0)

	return r.Listar(ctx, FiltroLancamento{
		Status:               "pending",
		Tipo:                 "payable",
		DataVencimentoInicio: &inicio,
		DataVencimentoFim:    &fim,
	})
}

func (r *LancamentoRepository) ListarPorPeriodo(ctx context.Context, inicio, fim time.Time, tipo string) ([]model.Entry, error) {
	if strings.TrimSpace(tipo) == "" {
		return nil, fmt.Errorf("tipo do lancamento e obrigatorio")
	}

	return r.Listar(ctx, FiltroLancamento{
		Status:               "pending",
		Tipo:                 tipo,
		DataVencimentoInicio: &inicio,
		DataVencimentoFim:    &fim,
	})
}
