package address

import (
	"net/http"

	"github.com/google/uuid"
	service "github.com/matheusgosk8/book-me-server/internal/service/address"
	"github.com/matheusgosk8/book-me-server/internal/utils"
)

func ListMyAddresses(w http.ResponseWriter, r *http.Request) {
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

	addressService := service.NewAddressService()
	addresses, err := addressService.ListMyAddresses(r.Context(), service.ListInput{UserID: userID})
	if err != nil {
		http.Error(w, "Erro ao buscar endereços", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	utils.ServerResponse(w, addresses.Addresses)

}
