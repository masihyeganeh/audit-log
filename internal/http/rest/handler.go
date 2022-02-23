package rest

import (
	"github.com/masihyeganeh/audit-log/internal/service"
)

type Handler struct {
	service service.Service
}

// CreateHandler Creates a new instance of REST handler
func CreateHandler(auditLogService service.Service) *Handler {
	return &Handler{
		service: auditLogService,
	}
}
