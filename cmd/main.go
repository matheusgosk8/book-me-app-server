package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	api "github.com/matheusgosk8/book-me-server/internal/router"
	log "github.com/sirupsen/logrus"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Error("Não foi possível carregar o arquivo .env, usando variáveis padrão")
	}

	log.SetReportCaller(true)
	var router *chi.Mux = chi.NewRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	api.Router(router)

	fmt.Printf("Starting book-me server at port: %s - ", port)

	err = http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	if err != nil {
		log.Error(err)
		log.Fatal("Server exits status 500")
	}
}
