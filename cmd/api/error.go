package main

import (
	"net/http"
	"task-api/internal/gateway"
	"task-api/pkg/webservice"
)

func mapError(code webservice.ErrCode, err error) (any, int) {
	switch code {
	case webservice.ErrCodeJsonParsing:
		return wrapMessage(err.Error()), http.StatusBadRequest
	case webservice.ErrCodeJsonBodyValidation:
		return wrapMessage(err.Error()), http.StatusBadRequest
	case webservice.ErrCodeMalformedEndpointHeader, webservice.ErrCodeUnsupportedEndpoint:
		return wrapMessage(err.Error()), http.StatusBadRequest
	case webservice.ErrCodeClientCode:
		if gatErr, ok := err.(*gateway.Error); ok {
			switch gatErr.Code() {
			case gateway.ErrCodeBadInput:
				return wrapMessage(err.Error()), http.StatusBadRequest
			case gateway.ErrCodeNotFound:
				return wrapMessage(err.Error()), http.StatusNotFound
			}
		}
	}
	return wrapMessage("что-то пошло не так"), http.StatusInternalServerError
}

func wrapMessage(msg string) map[string]any {
	return map[string]any{
		"error": msg,
	}
}
