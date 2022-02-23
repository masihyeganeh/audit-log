package rest

import (
	"net/http"
)

func (h *Handler) HealthCheck(req *http.Request, body []byte) (int, interface{}, error) {
	return http.StatusOK, "", nil
}
