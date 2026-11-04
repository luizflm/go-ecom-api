package products

import (
	"errors"
	"log"
	"net/http"
	"strconv"

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

func (h *handler) ListProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.service.ListProducts(r.Context())
	if err != nil {
		log.Println(err)
		json.WriteError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	json.Write(w, http.StatusOK, products)
}

func (h *handler) FindProductByID(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		json.WriteError(w, http.StatusUnprocessableEntity, "invalid ID")
		return
	}

	product, err := h.service.FindProductByID(r.Context(), id)
	if err != nil {
		switch {
		case errors.Is(err, ErrNotFound):
			json.WriteError(w, http.StatusNotFound, ErrNotFound.Error())
		default:
			log.Println(err)
			json.WriteError(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	json.Write(w, http.StatusOK, product)
}
