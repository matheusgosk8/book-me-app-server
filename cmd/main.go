package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	api "github.com/matheusgosk8/book-me-server/internal/router"
	log "github.com/sirupsen/logrus"
)

func main() {

	log.SetReportCaller(true)
	var router *chi.Mux = chi.NewRouter()
	const PORT int16 = 8000
	api.Router(router)

	fmt.Printf("Starting book-me server at port: %d - ", PORT)

	err := http.ListenAndServe(fmt.Sprintf(":%d", PORT), router)
	if err != nil {
		log.Error(err)
		log.Fatal("Server exits status 500")
	}
}
