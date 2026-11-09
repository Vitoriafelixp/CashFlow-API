package handler

type AtualizarOrcamentoReposicaoRequest struct {
	OrcamentoReposicao *float64 `json:"orcamento_reposicao"`
}

type OrcamentoReposicaoResponse struct {
	OrcamentoReposicao float64 `json:"orcamento_reposicao"`
}