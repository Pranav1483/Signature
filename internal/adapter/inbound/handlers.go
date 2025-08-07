package inbound

import "signature/internal/service"

type Handler struct {
	service *service.Service
	worker  chan func()
	pool    chan any
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		service: service,
		worker:  make(chan func(), 10),
		pool:    make(chan any, 10),
	}
}
