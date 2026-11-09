package handler

type SaldoLiquidoResponse struct {
	Dias         int     `json:"dias"`
	SaldoLiquido float64 `json:"saldo_liquido"`
}

type FaturamentoEsperadoResponse struct {
	Dias                int     `json:"dias"`
	FaturamentoEsperado float64 `json:"faturamento_esperado"`
}
