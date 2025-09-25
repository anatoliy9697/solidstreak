package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	apperrors "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/common/errors"
	habitPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/habit"
)

type GetHabitResponse struct {
	Data   *habitPkg.Habit   `json:"data,omitempty"`
	Errors []apperrors.Error `json:"errors,omitempty"`
}

func (s Server) getHabit(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var (
		err      error
		response GetHabitResponse
	)

	logger := s.Res.Logger

	// Adding request ID to request context
	reqID, _ := r.Context().Value(ctxKeyRequestID{}).(string)
	if reqID != "" {
		logger = logger.With("requestId", reqID)
	}

	defer func() {
		if err != nil {
			apperror, ok := err.(apperrors.Error)
			if !ok {
				apperror = apperrors.ErrInternal(err.Error())
			}

			logger.Error("error while getting habit", "error", err)

			w.WriteHeader(apperror.HTTPCode)
			response = GetHabitResponse{Errors: []apperrors.Error{apperror}}
			json.NewEncoder(w).Encode(response)
		}
	}()

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		err = apperrors.ErrBadRequest("missing habit id")
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		err = apperrors.ErrBadRequest("invalid habit id")
		return
	}

	habit, err := s.Res.HabitRepo.GetByID(id)
	if err != nil {
		return
	}

	response = GetHabitResponse{Data: habit}
	json.NewEncoder(w).Encode(response)
}
