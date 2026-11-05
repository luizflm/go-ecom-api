package orders

import (
	"errors"
	"log"
	"net/http"

	"github.com/luizflm/go-ecom-api/internal/json"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}
func (h *handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var tempOrder createOrderParams
	if err := json.Read(r, &tempOrder); err != nil {
		log.Println(err)
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validateCreateOrder(tempOrder); err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	createdOrder, err := h.service.PlaceOrder(r.Context(), tempOrder)
	if err != nil {
		log.Println(err)
		switch {
		case errors.Is(err, ErrProductNotFound):
			json.WriteError(w, http.StatusNotFound, ErrProductNotFound.Error())
		case errors.Is(err, ErrProductNoStock):
			json.WriteError(w, http.StatusConflict, ErrProductNoStock.Error())
		default:
			json.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	json.Write(w, http.StatusCreated, createdOrder)
}
