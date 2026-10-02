package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"cashflow-manager/internal/middleware"
	"cashflow-manager/internal/model"
	"cashflow-manager/internal/service"

	"github.com/go-chi/chi/v5"
)

type LancamentosResponse []model.Entry

type RelatorioHandler struct {
	service *service.RelatorioService
}

func NewRelatorioHandler(service *service.RelatorioService) *RelatorioHandler {
	return &RelatorioHandler{service: service}
}

func (h *RelatorioHandler) Registrar(router chi.Router) {
	router.Route("/relatorios", func(r chi.Router) {
		r.Get("/vencidos", h.ListarVencidos)
		r.Get("/saldo-liquido", h.SaldoLiquido)
		r.Get("/faturamento-esperado", h.FaturamentoEsperado)
	})
}

// ListarVencidos godoc
// @Summary      Listar lançamentos vencidos
// @Tags         relatorios
// @Produce      json
// @Success      200 {array} LancamentosResponse
// @Failure      500 {string} string
// @Router       /relatorios/vencidos [get]
func (h *RelatorioHandler) ListarVencidos(w http.ResponseWriter, r *http.Request) {
	lancamentos, err := h.service.ListarVencidos(r.Context())
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "relatorio_error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, LancamentosResponse(lancamentos))
}

// SaldoLiquido godoc
// @Summary      Obter saldo líquido
// @Tags         relatorios
// @Produce      json
// @Param        dias query int true "Quantidade de dias do horizonte"
// @Success      200 {object} SaldoLiquidoResponse
// @Failure      400 {string} string
// @Failure      500 {string} string
// @Router       /relatorios/saldo-liquido [get]
func (h *RelatorioHandler) SaldoLiquido(w http.ResponseWriter, r *http.Request) {
	dias, err := parseDias(r)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid_query", err.Error())
		return
	}

	saldo, err := h.service.SaldoLiquido(r.Context(), dias)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "relatorio_error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, SaldoLiquidoResponse{
		Dias:         dias,
		SaldoLiquido: saldo,
	})
}

// FaturamentoEsperado godoc
// @Summary      Obter faturamento esperado
// @Tags         relatorios
// @Produce      json
// @Param        dias query int true "Quantidade de dias do horizonte"
// @Success      200 {object} FaturamentoEsperadoResponse
// @Failure      400 {string} string
// @Failure      500 {string} string
// @Router       /relatorios/faturamento-esperado [get]
func (h *RelatorioHandler) FaturamentoEsperado(w http.ResponseWriter, r *http.Request) {
	dias, err := parseDias(r)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid_query", err.Error())
		return
	}

	valor, err := h.service.FaturamentoEsperado(r.Context(), dias)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "relatorio_error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, FaturamentoEsperadoResponse{
		Dias:                dias,
		FaturamentoEsperado: valor,
	})
}

func parseDias(r *http.Request) (int, error) {
	valor := r.URL.Query().Get("dias")
	if valor == "" {
		return 0, fmt.Errorf("parametro 'dias' e obrigatorio")
	}

	dias, err := strconv.Atoi(valor)
	if err != nil || dias < 0 {
		return 0, fmt.Errorf("parametro 'dias' deve ser um inteiro >= 0")
	}

	return dias, nil
}
