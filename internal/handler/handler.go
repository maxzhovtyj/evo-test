package handler

import (
	"evo-test/internal/service"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type Handler interface {
	Register(router *httprouter.Router)
}

type handler struct {
	service service.Service
}

func New(service service.Service) Handler {
	return &handler{service: service}
}

func (h *handler) Register(router *httprouter.Router) {
	router.HandlerFunc(http.MethodGet, "/transaction", h.GetTransaction)
}

func (h *handler) GetTransaction(w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}
