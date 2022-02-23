package rest

import (
	"encoding/json"
	"github.com/masihyeganeh/audit-log/internal/auth"
	"github.com/masihyeganeh/audit-log/internal/service"
	"github.com/pkg/errors"
	"net/http"
)

func (h *Handler) Log(req *http.Request, body []byte) (int, interface{}, error) {
	u, ok := req.Context().Value("auth").(*auth.User)
	if !ok {
		return http.StatusForbidden, nil, errors.New("forbidden")
	}

	if !u.HasWriteAccess {
		return http.StatusForbidden, nil, errors.New("you don't have write access")
	}

	var logRequest LogRequest
	err := json.Unmarshal(body, &logRequest)
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrap(err, "could not parse log request")
	}

	err = h.service.LogEvent(req.Context(), service.Event{
		EventType:    logRequest.EventType,
		CommonField1: logRequest.CommonField1,
		CommonField2: logRequest.CommonField2,
		Fields:       logRequest.Fields,
	})
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	return http.StatusOK, "ok", nil
}
