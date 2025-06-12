package webservice

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
)

type ErrCode int

const (
	ErrCodeJsonParsing = iota
	ErrCodeJsonBodyValidation
	ErrCodeMalformedEndpointHeader
	ErrCodeUnsupportedEndpoint
	ErrCodeClientCode
)

type ErrorMapper = func(code ErrCode, e error) (any, int)

type Validator interface {
	Validate() error
}

type service struct {
	handlers  map[string]http.HandlerFunc
	errMapper ErrorMapper
}

func New() *service {
	return &service{
		handlers: make(map[string]http.HandlerFunc),
	}
}

func (s *service) WithErrorMapper(m ErrorMapper) {
	s.errMapper = m
}

// Registers handler on an endpoint
func Register[T any, U any](s *service, endpoint string, h func(context.Context, *T, *U) error) {
	s.handlers[endpoint] = func(w http.ResponseWriter, r *http.Request) {
		var req T
		if reflect.TypeFor[T]().NumField() > 0 {
			dec := json.NewDecoder(r.Body)
			err := dec.Decode(&req)
			if err != nil {
				s.writeError(w, ErrCodeJsonParsing, err)
				return
			}
			if v, ok := any(req).(Validator); ok {
				err := v.Validate()
				if err != nil {
					s.writeError(w, ErrCodeJsonBodyValidation, err)
					return
				}
			}
		}
		var res U
		err := h(context.Background(), &req, &res)
		if err != nil {
			s.writeError(w, ErrCodeClientCode, err)
			return
		}
		enc := json.NewEncoder(w)
		enc.Encode(res)
	}
	logMsg := fmt.Sprintf("%s эндпоинт зарегистрирован (%s -> %s).", endpoint, reflect.TypeFor[T](), reflect.TypeFor[U]())
	slog.Info(logMsg)
}

func (s *service) endpoints() []string {
	endpoints := make([]string, 0, len(s.handlers))
	for endpoint := range s.handlers {
		endpoints = append(endpoints, endpoint)
	}
	return endpoints
}

func (s *service) Handle(w http.ResponseWriter, r *http.Request) {
	endpoint := r.Header["Endpoint"]
	w.Header().Set("Content-Type", "application/json")
	if len(endpoint) != 1 {
		err := fmt.Errorf(
			"заголовок `Endpoint` запроса должен содержать название эндпоинта. Поддерживаемые эндпоинты: %s", s.endpoints(),
		)
		s.writeError(w, ErrCodeMalformedEndpointHeader, err)
		return
	}
	handler, ok := s.handlers[endpoint[0]]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		err := fmt.Errorf(
			"эндпоинт `%s` не поддерживается. Поддерживаемые эндпоинты: %s", endpoint[0], s.endpoints(),
		)
		s.writeError(w, ErrCodeUnsupportedEndpoint, err)
		return
	}
	handler(w, r)
}

func (s *service) writeError(w http.ResponseWriter, errCode ErrCode, err error) {
	enc := json.NewEncoder(w)
	if s.errMapper != nil {
		rsp, code := s.errMapper(errCode, err)
		w.WriteHeader(code)
		enc.Encode(rsp)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		enc.Encode(map[string]any{
			"code":    errCode,
			"message": err.Error(),
		})
	}
}
