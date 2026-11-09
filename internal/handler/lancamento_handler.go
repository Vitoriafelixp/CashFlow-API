package handler

import (
	"errors"
	"net/http"

	"cashflow-manager/internal/middleware"
	"cashflow-manager/internal/service"

	"github.com/go-chi/chi/v5"
)

type LancamentoHandler struct {
	service *service.LancamentoService
}

func NewLancamentoHandler(service *service.LancamentoService) *LancamentoHandler {
	return &LancamentoHandler{service: service}
}

func (h *LancamentoHandler) Registrar(router chi.Router) {
	router.Get("/lancamentos", h.Listar)
	router.Post("/lancamentos", h.Criar)
}

// Listar godoc
// @Summary      Listar lançamentos
// @Tags         lancamentos
// @Produce      json
// @Success      200 {array} ListaLancamentosResponse
// @Failure      500 {string} string
// @Router       /lancamentos [get]
func (h *LancamentoHandler) Listar(w http.ResponseWriter, r *http.Request) {
	lancamentos, err := h.service.Listar(r.Context())
	if err != nil {
		middleware.WriteError(w, http.StatusInternalServerError, "lancamento_error", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, toListaLancamentosResponse(lancamentos))
}

// Criar godoc
// @Summary      Criar lançamento
// @Tags         lancamentos
// @Accept       json
// @Produce      json
// @Param        request body CriarLancamentoRequest true "Dados do lançamento"
// @Success      201 {object} LancamentoResponse
// @Failure      400 {string} string
// @Failure      500 {string} string
// @Router       /lancamentos [post]
func (h *LancamentoHandler) Criar(w http.ResponseWriter, r *http.Request) {
	var payload CriarLancamentoRequest
	if err := decodeJSON(r, &payload); err != nil {
		middleware.WriteError(w, http.StatusBadRequest, "invalid_json", "Request body must contain valid JSON.")
		return
	}

	lancamento, err := h.service.Criar(r.Context(), payload.toModel())
	if err != nil {
		if errors.Is(err, service.ErrInvalidLancamento) {
			middleware.WriteError(w, http.StatusBadRequest, "invalid_payload", err.Error())
			return
		}
		middleware.WriteError(w, http.StatusInternalServerError, "lancamento_error", err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, toLancamentoResponse(lancamento))
}
