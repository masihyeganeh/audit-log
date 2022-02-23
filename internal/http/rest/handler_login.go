package rest

import (
	"encoding/json"
	"github.com/masihyeganeh/audit-log/internal/service"
	"github.com/pkg/errors"
	"net/http"
)

func (h *Handler) Login(req *http.Request, body []byte) (int, interface{}, error) {
	var loginRequest LoginRequest
	err := json.Unmarshal(body, &loginRequest)
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrap(err, "could not parse login request")
	}

	jwtToken, err := h.service.Login(req.Context(), service.LoginRequest(loginRequest))
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, map[string]string{"jwt_token": jwtToken}, nil
}
