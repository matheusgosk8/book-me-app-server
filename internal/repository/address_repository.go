package repository

import (
	"context"

	"github.com/matheusgosk8/book-me-server/ent"
	api "github.com/matheusgosk8/book-me-server/internal/models"
	log "github.com/sirupsen/logrus"
)

func CreateAddress(ctx context.Context, client *ent.Client, params api.Address) (*ent.Address, error) {
	u, err := client.Address.
		Create().
		SetLabel(params.Street).
		SetCity(params.City).
		SetState(params.State).
		SetPostalCode(params.PostalCode).
		SetCountry(params.Country).
		Save(ctx)
	if err != nil {
		log.Errorf("Error creating address: %v", err)
		return nil, err
	}
	return u, nil
}

// Id           string    `json:"id"`
// Street       string    `json:"street"`
// City         string    `json:"city"`
// State        string    `json:"state"`
// PostalCode   string    `json:"postal_code"`
// Country      string    `json:"country"`
// CreationDate time.Time `json:"creation_date"`
