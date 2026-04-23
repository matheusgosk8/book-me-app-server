package models

import "time"

type Address struct {
	Id           string    `json:"id"`
	Street       string    `json:"street"`
	City         string    `json:"city"`
	State        string    `json:"state"`
	PostalCode   string    `json:"postal_code"`
	Country      string    `json:"country"`
	CreationDate time.Time `json:"creation_date"`
}
