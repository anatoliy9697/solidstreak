package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	apperrors "github.com/anatoliy9697/solidstreak/solidstreak-backend/pkg/errors"

	hPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/habit"
	hRepo "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/habit/repo"
	usrPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/user"
)

type GetHabitResponse struct {
	Data *hPkg.Habit `json:"data"`
}

type Habit struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
}

type PostPutHabitRequest struct {
	Data Habit    `json:"data"`
	Meta Metadata `json:"meta"`
}

type PostPutHabitResponse struct {
	Data *hPkg.Habit `json:"data"`
}

type GetHabitsResponse struct {
	Data []*hPkg.Habit `json:"data"`
}

// TODO: попробовать избавиться от дублирования кода

func (s Server) getHabit(w http.ResponseWriter, r *http.Request) {
	var err error

	logger := s.Res.Logger

	// Adding request ID to request context
	reqID, _ := r.Context().Value(ctxKeyRequestID{}).(string)
	if reqID != "" {
		logger = logger.With("requestId", reqID)
	}

	defer func() {
		if err != nil {
			processError(w, logger, err)
		}
	}()

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		err = apperrors.ErrBadRequest("missing habit id")
		return
	}

	var id int64
	id, err = strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		err = apperrors.ErrBadRequest("invalid habit id")
		return
	}

	userTgID, ok := r.Context().Value(ctxKeyUserTgID{}).(int64)
	if !ok {
		err = apperrors.ErrUnauthorized("couldn't identify user")
		return
	}

	var user *usrPkg.User
	if user, err = s.Res.UsrRepo.GetByTgID(userTgID); err != nil {
		return
	}

	var habit *hPkg.Habit
	if habit, err = s.Res.HabitRepo.GetByIDAndOwnerID(id, user.ID); err != nil {
		return
	}

	response := GetHabitResponse{Data: habit}

	json.NewEncoder(w).Encode(response)
}

func (s Server) postHabit(w http.ResponseWriter, r *http.Request) {
	var err error

	logger := s.Res.Logger

	// Adding request ID to request context
	reqID, _ := r.Context().Value(ctxKeyRequestID{}).(string)
	if reqID != "" {
		logger = logger.With("requestId", reqID)
	}

	defer func() {
		if err != nil {
			processError(w, logger, err)
		}
	}()

	var req PostPutHabitRequest

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&req); err != nil {
		err = apperrors.ErrBadRequest("invalid request payload")
		return
	}
	if req.Data.Title == "" {
		err = apperrors.ErrBadRequest("habit title is required")
		return
	}
	if req.Meta.UserID == 0 {
		err = apperrors.ErrBadRequest("userId in meta is required")
		return
	}

	userTgID, ok := r.Context().Value(ctxKeyUserTgID{}).(int64)
	if !ok {
		err = apperrors.ErrUnauthorized("couldn't identify user")
		return
	}

	var user *usrPkg.User
	if user, err = s.Res.UsrRepo.GetByTgID(userTgID); err != nil {
		return
	}

	if user.ID != req.Meta.UserID {
		err = apperrors.ErrForbidden("userId in meta does not match the authenticated user")
		return
	}

	habit := hPkg.NewHabit(req.Data.Title, req.Data.Description, req.Meta.UserID)

	if err = s.Res.HabitRepo.Create(habit); err != nil {
		return
	}

	response := PostPutHabitResponse{Data: habit}

	json.NewEncoder(w).Encode(response)
}

func (s Server) putHabit(w http.ResponseWriter, r *http.Request) {
	var err error

	logger := s.Res.Logger

	// Adding request ID to request context
	reqID, _ := r.Context().Value(ctxKeyRequestID{}).(string)
	if reqID != "" {
		logger = logger.With("requestId", reqID)
	}

	defer func() {
		if err != nil {
			processError(w, logger, err)
		}
	}()

	idStr := chi.URLParam(r, "id")
	if idStr == "" {
		err = apperrors.ErrBadRequest("missing habit id")
		return
	}

	var id int64
	id, err = strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		err = apperrors.ErrBadRequest("invalid habit id")
		return
	}

	var req PostPutHabitRequest

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&req); err != nil {
		err = apperrors.ErrBadRequest("invalid request payload")
		return
	}
	if req.Data.Title == "" {
		err = apperrors.ErrBadRequest("habit title is required")
		return
	}
	if req.Meta.UserID == 0 {
		err = apperrors.ErrBadRequest("userId in meta is required")
		return
	}

	userTgID, ok := r.Context().Value(ctxKeyUserTgID{}).(int64)
	if !ok {
		err = apperrors.ErrUnauthorized("couldn't identify user")
		return
	}

	var user *usrPkg.User
	if user, err = s.Res.UsrRepo.GetByTgID(userTgID); err != nil {
		return
	}

	if user.ID != req.Meta.UserID {
		err = apperrors.ErrForbidden("userId in meta does not match the authenticated user")
		return
	}

	habit, err := s.Res.HabitRepo.GetByIDAndOwnerID(id, req.Meta.UserID)
	if err != nil {
		return
	}

	habit.Title = req.Data.Title
	habit.Description = req.Data.Description
	habit.UpdatedAt = time.Now()

	if err = s.Res.HabitRepo.Update(habit); err != nil {
		return
	}

	response := PostPutHabitResponse{Data: habit}

	json.NewEncoder(w).Encode(response)
}

func (s Server) getHabits(w http.ResponseWriter, r *http.Request) {
	var err error

	logger := s.Res.Logger

	// Adding request ID to request context
	reqID, _ := r.Context().Value(ctxKeyRequestID{}).(string)
	if reqID != "" {
		logger = logger.With("requestId", reqID)
	}

	defer func() {
		if err != nil {
			processError(w, logger, err)
		}
	}()

	userTgID, ok := r.Context().Value(ctxKeyUserTgID{}).(int64)
	if !ok {
		err = apperrors.ErrUnauthorized("couldn't identify user")
		return
	}

	var user *usrPkg.User
	if user, err = s.Res.UsrRepo.GetByTgID(userTgID); err != nil {
		return
	}

	var habits []*hPkg.Habit
	if habits, err = s.Res.HabitRepo.GetByOwnerID(user.ID, hRepo.Any); err != nil { // TODO: получать gettingMode из query параметра
		return
	}

	response := GetHabitsResponse{Data: habits}

	json.NewEncoder(w).Encode(response)
}
