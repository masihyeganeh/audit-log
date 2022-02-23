package rest

import (
	"encoding/json"
	"github.com/masihyeganeh/audit-log/internal/auth"
	"github.com/pkg/errors"
	"net/http"
)

func (h *Handler) Query(req *http.Request, body []byte) (int, interface{}, error) {
	u, ok := req.Context().Value("auth").(*auth.User)
	if !ok {
		return http.StatusForbidden, nil, errors.New("forbidden")
	}

	if !u.HasReadAccess {
		return http.StatusForbidden, nil, errors.New("you don't have read access")
	}

	var queryRequest QueryRequest
	err := json.Unmarshal(body, &queryRequest)
	if err != nil {
		return http.StatusBadRequest, nil, errors.Wrap(err, "could not parse log request")
	}

	r, err := h.service.QueryEvents(req.Context(), queryRequest.EventType, queryRequest.Filters)
	if err != nil {
		return http.StatusInternalServerError, nil, err
	}

	events := make([]LogResponse, len(r))
	for i, event := range r {
		events[i] = LogResponse{
			EventTime:    event.EventTime,
			EventType:    event.EventType,
			CommonField1: event.CommonField1,
			CommonField2: event.CommonField2,
			Fields:       event.Fields,
		}
	}

	return http.StatusOK, r, nil
}
