package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/masihyeganeh/audit-log/internal/http/rest"
	"io/ioutil"
	"net/http"
	"strings"
)

// SetupRouter creates routes
func (s *server) SetupRouter(ctx context.Context, handler *rest.Handler) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/login", s.WrapHandler(ctx, handler.Login))

	mux.HandleFunc("/log", s.WrapHandler(ctx, handler.Log))

	mux.HandleFunc("/query", s.WrapHandler(ctx, handler.Query))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("Audit-log project"))
	})

	mux.HandleFunc("/health-check", s.WrapHandler(ctx, handler.HealthCheck))

	return mux
}

type HandlerFunc func(req *http.Request, body []byte) (int, interface{}, error)

type Response struct {
	Status   string      `json:"status"`
	Response interface{} `json:"response"`
	Error    string      `json:"error"`
}

func (s *server) WrapHandler(ctx context.Context, handler HandlerFunc) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Read body of request
		var body []byte
		if r.Body != nil {
			var err error
			body, err = ioutil.ReadAll(r.Body)
			if err != nil {
				body = []byte{}
			}
		}

		// Authenticate
		authorization := r.Header.Get("Authorization")
		if len(authorization) > 0 && strings.HasPrefix(authorization, "Bearer ") {
			authorization = strings.Replace(authorization, "Bearer ", "", 1)
			if user, err := s.authentication.Authenticate(authorization); err == nil {
				ctx = context.WithValue(ctx, "auth", user)
			}
		}

		r = r.WithContext(ctx)

		// Call the actual handler
		statusCode, response, err := handler(r, body)

		res := Response{
			Status:   "",
			Response: nil,
			Error:    "",
		}

		if err != nil {
			res.Status = "error"
			res.Error = err.Error()
		} else {
			res.Status = "ok"
			res.Response = response
		}

		body, err = json.Marshal(res)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(fmt.Sprintf("could not marshal response : %v", err)))
			return
		}

		w.WriteHeader(statusCode)
		_, _ = w.Write(body)
	}
}
