package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	apperrors "github.com/anatoliy9697/solidstreak/solidstreak-backend/pkg/errors"
)

type ErrorResponse struct {
	Errors []*apperrors.Error `json:"errors"`
}

func writeError(w http.ResponseWriter, err error) {
	var apperror *apperrors.Error
	if !errors.As(err, &apperror) {
		apperror = apperrors.NewInternalErr(err.Error())
	}

	w.WriteHeader(apperror.HTTPCode)
	response := ErrorResponse{Errors: []*apperrors.Error{apperror}}
	json.NewEncoder(w).Encode(response)
}

func processError(w http.ResponseWriter, logger *slog.Logger, err error) {
	var apperror *apperrors.Error
	if !errors.As(err, &apperror) {
		apperror = apperrors.NewInternalErr(err.Error())
	}

	logger.Error("error occurred", "error", err)

	writeError(w, apperror)
}
