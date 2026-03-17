package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/matheusgosk8/book-me-server/internal/models"
	"github.com/matheusgosk8/book-me-server/internal/utils"

	log "github.com/sirupsen/logrus"
)

func Status(res http.ResponseWriter, req *http.Request) {

	start := time.Now()

	log.Infof("Accessed route: %s - %s | Duration: %s", req.Method, req.URL.Path, time.Since(start))

	var response = models.ServerStatus{
		Code:    200,
		Message: "Server health ok!",
	}

	res.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(res).Encode(response)

	log.Info("Responding request with ", response)

	if err != nil {
		utils.InternalErrorHandler(res, err)
	}

}
