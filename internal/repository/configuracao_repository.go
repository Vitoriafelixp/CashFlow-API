package repository

import (
	"context"
	"database/sql"
	"fmt"

	"cashflow-manager/internal/database"
	"cashflow-manager/internal/model"
)

type ConfiguracaoRepository struct {
	db *database.DB
}

func NewConfiguracaoRepository(db *database.DB) *ConfiguracaoRepository {
	return &ConfiguracaoRepository{db: db}
}

func (r *ConfiguracaoRepository) ObterOrcamentoReposicao(ctx context.Context) (float64, error) {
	var configuracao model.Settings

	err := r.db.GetContext(ctx, &configuracao, `
		SELECT id, restock_budget
		FROM settings
		ORDER BY id DESC
		LIMIT 1
	`)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil
		}
		return 0, fmt.Errorf("obter orcamento reposicao: %w", err)
	}

	return configuracao.RestockBudget, nil
}

func (r *ConfiguracaoRepository) AtualizarOrcamentoReposicao(ctx context.Context, valor float64) (float64, error) {
	var configuracao model.Settings

	err := r.db.GetContext(ctx, &configuracao, `
		SELECT id, restock_budget
		FROM settings
		ORDER BY id DESC
		LIMIT 1
	`)
	if err != nil {
		if err == sql.ErrNoRows {
			result, insertErr := r.db.ExecContext(ctx, `
				INSERT INTO settings (restock_budget)
				VALUES (?)
			`, valor)
			if insertErr != nil {
				return 0, fmt.Errorf("inserir orcamento reposicao: %w", insertErr)
			}

			if _, err := result.LastInsertId(); err != nil {
				return 0, fmt.Errorf("ler id do orcamento reposicao inserido: %w", err)
			}

			return valor, nil
		}
		return 0, fmt.Errorf("consultar configuracao existente: %w", err)
	}

	_, err = r.db.ExecContext(ctx, `
		UPDATE settings
		SET restock_budget = ?
		WHERE id = ?
	`, valor, configuracao.ID)
	if err != nil {
		return 0, fmt.Errorf("atualizar orcamento reposicao: %w", err)
	}

	return valor, nil
}
