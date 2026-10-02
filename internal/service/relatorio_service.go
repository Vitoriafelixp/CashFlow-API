package service

import (
	"context"
	"fmt"
	"time"

	"cashflow-manager/internal/model"
	"cashflow-manager/internal/repository"
)

type RelatorioService struct {
	lancamentoRepository   *repository.LancamentoRepository
	configuracaoRepository *repository.ConfiguracaoRepository
}

func NewRelatorioService(lancamentoRepository *repository.LancamentoRepository, configuracaoRepository *repository.ConfiguracaoRepository) *RelatorioService {
	return &RelatorioService{
		lancamentoRepository:   lancamentoRepository,
		configuracaoRepository: configuracaoRepository,
	}
}

func (s *RelatorioService) SaldoLiquido(ctx context.Context, dias int) (float64, error) {
	if s == nil || s.lancamentoRepository == nil {
		return 0, fmt.Errorf("lancamento repository nao inicializado")
	}
	if dias < 0 {
		return 0, fmt.Errorf("dias deve ser maior ou igual a zero")
	}

	inicio := time.Now().Truncate(24 * time.Hour)
	fim := inicio.AddDate(0, 0, dias)

	receitas, err := s.somarPorFiltro(ctx, repository.FiltroLancamento{
		Status:                        "pending",
		Tipo:                          "receivable",
		DataPrevistaRecebimentoInicio: &inicio,
		DataPrevistaRecebimentoFim:    &fim,
	})
	if err != nil {
		return 0, err
	}

	despesas, err := s.somarPorFiltro(ctx, repository.FiltroLancamento{
		Status:               "pending",
		Tipo:                 "payable",
		DataVencimentoInicio: &inicio,
		DataVencimentoFim:    &fim,
	})
	if err != nil {
		return 0, err
	}

	return receitas - despesas, nil
}

func (s *RelatorioService) FaturamentoEsperado(ctx context.Context, dias int) (float64, error) {
	if s == nil || s.lancamentoRepository == nil || s.configuracaoRepository == nil {
		return 0, fmt.Errorf("repositories de relatorio nao inicializados")
	}
	if dias < 0 {
		return 0, fmt.Errorf("dias deve ser maior ou igual a zero")
	}

	inicio := time.Now().Truncate(24 * time.Hour)
	fim := inicio.AddDate(0, 0, dias)

	orcamentoReposicao, err := s.configuracaoRepository.ObterOrcamentoReposicao(ctx)
	if err != nil {
		return 0, err
	}

	despesas, err := s.somarPorFiltro(ctx, repository.FiltroLancamento{
		Status:               "pending",
		Tipo:                 "payable",
		DataVencimentoInicio: &inicio,
		DataVencimentoFim:    &fim,
	})
	if err != nil {
		return 0, err
	}

	confirmadas, err := s.somarPorFiltro(ctx, repository.FiltroLancamento{
		Status:                        "pending",
		Tipo:                          "receivable",
		Confirmado:                    boolPtr(true),
		DataPrevistaRecebimentoInicio: &inicio,
		DataPrevistaRecebimentoFim:    &fim,
	})
	if err != nil {
		return 0, err
	}

	return despesas + orcamentoReposicao - confirmadas, nil
}

func (s *RelatorioService) somarPorFiltro(ctx context.Context, filtro repository.FiltroLancamento) (float64, error) {
	lancamentos, err := s.lancamentoRepository.Listar(ctx, filtro)
	if err != nil {
		return 0, err
	}

	var total float64
	for _, lancamento := range lancamentos {
		total += lancamento.Amount
	}

	return total, nil
}

func boolPtr(v bool) *bool {
	return &v
}

func (s *RelatorioService) ListarVencidos(ctx context.Context) ([]model.Entry, error) {
	if s == nil || s.lancamentoRepository == nil {
		return nil, fmt.Errorf("lancamento repository nao inicializado")
	}

	return s.lancamentoRepository.Listar(ctx, repository.FiltroLancamento{
		Status:               "pending",
		Tipo:                 "payable",
		DataVencimentoInicio: datePtr(time.Now().AddDate(0, 0, -30)),
		DataVencimentoFim:    datePtr(time.Now()),
	})
}

func datePtr(t time.Time) *time.Time {
	return &t
}
