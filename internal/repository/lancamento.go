package repository

import (
	"context"
	"github.com/jmoiron/sqlx"
	"database/sql"
	"errors"
	"internal/model"
)

type LancamentoRepository struct {
	db *sqlx.DB
}

func NewLancamentoRepository(db sqlx.DB) *LancamentoRepository {
	return &LancamentoRepository{db: &db}
}

func (r *LancamentoRepository) Criar( ctx context.Context, Lancamento *model.Lancamento) error {
	query := `
	INSERT INTO entries (description, amount, type, due_date, expected_receipt_date, category, status, paid_at) 
	VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.ExecContext(ctx, query,
		Lancamento.Descricao,
		Lancamento.Valor,
		Lancamento.Tipo,
		Lancamento.DataVencimento,
		Lancamento.DataPrevistaRecebimento,
		Lancamento.Categoria,
		Lancamento.Status,
		Lancamento.Confirmado,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	Lancamento.ID = id

	return nil
}

func (r *LancamentoRepository) BuscarPorID(ctx context.Context, id int64) (*model.Lancamento, error) {
	query := `
	SELECT id, description, amount, type, due_date, expected_receipt_date, category, status, paid_at, confirmed, created_at
	FROM entries
	WHERE id = ?
	`

	var lancamento model.Lancamento
	err := r.db.GetContext(ctx, &lancamento, query, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("lançamento não encontrado")
		}
		return nil, err
	}

	return &lancamento, nil

}