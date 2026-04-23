package config

import (
	"github.com/go-chi/cors"
)

func GetCORSConfig() cors.Options{
	return cors.Options{
		AllowedOrigins: []string{
			"http://localhost:8081",
			"http://192.168.1.108:8081",
			"http://192.168.1.108:8000",
		},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{
		"Accept",
		"Content-Type",
		"Content-Lenght",
		"Content-Encoding",
		"Authorization",
		"x-CSRF-Token",
		},
		AllowCredentials: true,
		MaxAge: 300, 
	}
}