package service

import (
	"context"
	"fmt"

	"cashflow-manager/internal/repository"
)

type ConfiguracaoService struct {
	repository *repository.ConfiguracaoRepository
}

func NewConfiguracaoService(repository *repository.ConfiguracaoRepository) *ConfiguracaoService {
	return &ConfiguracaoService{repository: repository}
}

func (s *ConfiguracaoService) ObterOrcamentoReposicao(ctx context.Context) (float64, error) {
	if s == nil || s.repository == nil {
		return 0, fmt.Errorf("configuracao repository nao inicializado")
	}

	return s.repository.ObterOrcamentoReposicao(ctx)
}

func (s *ConfiguracaoService) AtualizarOrcamentoReposicao(ctx context.Context, valor float64) (float64, error) {
	if s == nil || s.repository == nil {
		return 0, fmt.Errorf("configuracao repository nao inicializado")
	}

	if valor < 0 {
		return 0, fmt.Errorf("orcamento de reposicao nao pode ser negativo")
	}

	return s.repository.AtualizarOrcamentoReposicao(ctx, valor)
}
