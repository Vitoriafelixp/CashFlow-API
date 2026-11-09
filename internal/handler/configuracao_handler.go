package handler

import (
	"net/http"

	"cashflow-manager/internal/middleware"
	"cashflow-manager/internal/service"

	"github.com/go-chi/chi/v5"
)

type ConfiguracaoHandler struct {
	service *service.ConfiguracaoService
}

func NewConfiguracaoHandler(service *service.ConfiguracaoService) *ConfiguracaoHandler {
	return &ConfiguracaoHandler{service: service}
}

func (h *ConfiguracaoHandler) Registrar(router chi.Router) {
	router.Route("/configuracoes", func(r chi.Router) {
		r.Get("/orcamento-reposicao", h.ObterOrcamentoReposicao)
		r.Put("/orcamento-reposicao", h.AtualizarOrcamentoReposicao)
	})
}

// ObterOrcamentoReposicao godoc
// @Summary      Obter orçamento de reposição
// @Tags         configuracoes
// @Produce      json
// @Success      200 {object} OrcamentoReposicaoResponse
// @Failure      500 {string} string
// @Router       /configuracoes/orcamento-reposicao [get]
func (h *ConfiguracaoHandler) ObterOrcamentoReposicao(w http.ResponseWriter, r *http.Request) {
	valor, err := h.service.ObterOrcamentoReposicao(r.Context())
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "configuracao_error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, OrcamentoReposicaoResponse{OrcamentoReposicao: valor})
}

// AtualizarOrcamentoReposicao godoc
// @Summary      Atualizar orçamento de reposição
// @Tags         configuracoes
// @Accept       json
// @Produce      json
// @Param        request body AtualizarOrcamentoReposicaoRequest true "Payload"
// @Success      200 {object} OrcamentoReposicaoResponse
// @Failure      400 {string} string
// @Failure      500 {string} string
// @Router       /configuracoes/orcamento-reposicao [put]
func (h *ConfiguracaoHandler) AtualizarOrcamentoReposicao(w http.ResponseWriter, r *http.Request) {
	var payload AtualizarOrcamentoReposicaoRequest
	if err := decodeJSON(r, &payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid_json", "Request body must contain valid JSON.")
		return
	}

	if payload.OrcamentoReposicao == nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid_payload", "Field orcamento_reposicao is required.")
		return
	}

	valor, err := h.service.AtualizarOrcamentoReposicao(r.Context(), *payload.OrcamentoReposicao)
	if err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid_payload", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, OrcamentoReposicaoResponse{OrcamentoReposicao: valor})
}
