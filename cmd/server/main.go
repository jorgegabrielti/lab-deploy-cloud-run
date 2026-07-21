package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jorgegabrielti/desafio-clima-cep/pkg/handler"
	"github.com/jorgegabrielti/desafio-clima-cep/pkg/viacep"
	"github.com/jorgegabrielti/desafio-clima-cep/pkg/weather"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	viaClient := viacep.NewClient()
	weaClient := weather.NewClient()
	weatherHandler := handler.NewWeatherHandler(viaClient, weaClient)

	mux := http.NewServeMux()
	mux.Handle("/", weatherHandler)
	mux.Handle("/weather", weatherHandler)
	mux.Handle("/weather/", weatherHandler)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	go func() {
		log.Printf("🚀 Servidor Clima por CEP iniciado na porta %s...", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar servidor HTTP: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Encerrando servidor graciosamente...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Erro no shutdown do servidor: %v", err)
	}

	log.Println("👋 Servidor finalizado com sucesso.")
}
