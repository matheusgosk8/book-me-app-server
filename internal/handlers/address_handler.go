package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/matheusgosk8/book-me-server/internal/db"
	"github.com/matheusgosk8/book-me-server/ent/address"
	"github.com/matheusgosk8/book-me-server/ent/user"
)

func ListMyAddresses(w http.ResponseWriter, r *http.Request) {
    // Recupera o ID do contexto. Como o Go é fortemente tipado, 
    // verificamos se ele já é um UUID ou se precisa de conversão.
    var userID uuid.UUID
    
    switch v := r.Context().Value("user_id").(type) {
    case uuid.UUID:
        userID = v
    case string:
        userID, _ = uuid.Parse(v)
    default:
        http.Error(w, "Não foi possível identificar o usuário", http.StatusUnauthorized)
        return
    }

    // Busca usando a Edge 'user' definida no schema de Address[cite: 4]
    addresses, err := db.Client.Address.
        Query().
        Where(address.HasUserWith(user.IDEQ(userID))).
        All(r.Context())

    if err != nil {
        http.Error(w, "Erro ao buscar endereços", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(addresses)
}