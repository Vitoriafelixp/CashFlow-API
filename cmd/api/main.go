// @title           CashFlow Manager API
// @version         1.0
// @description     API para gestão de fluxo de caixa com relatórios de saldo líquido e faturamento esperado.
// @host            localhost:8080
// @BasePath        /
package main

import (
	"log"
	"net/http"
	"os"

	"cashflow-manager/internal/database"
	"cashflow-manager/internal/handler"
	"cashflow-manager/internal/middleware"
	"cashflow-manager/internal/repository"
	"cashflow-manager/internal/service"

	"github.com/go-chi/chi/v5"
)

func main() {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Fatal("DATABASE_DSN must be set")
	}

	db, err := database.NewDatabase(dsn)
	if err != nil {
		log.Fatal(err)
	}

	configuracaoRepository := repository.NewConfiguracaoRepository(db)
	lancamentoRepository := repository.NewLancamentoRepository(db)

	configuracaoService := service.NewConfiguracaoService(configuracaoRepository)
	lancamentoService := service.NewLancamentoService(lancamentoRepository)
	relatorioService := service.NewRelatorioService(lancamentoRepository, configuracaoRepository)

	configuracaoHandler := handler.NewConfiguracaoHandler(configuracaoService)
	lancamentoHandler := handler.NewLancamentoHandler(lancamentoService)
	relatorioHandler := handler.NewRelatorioHandler(relatorioService)

	router := chi.NewRouter()
	router.Use(
		middleware.RequestLogger,
		middleware.Recoverer,
		middleware.ValidateJSON,
	)

	configuracaoHandler.Registrar(router)
	lancamentoHandler.Registrar(router)
	relatorioHandler.Registrar(router)

	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	log.Fatal(http.ListenAndServe(addr, router))
}
