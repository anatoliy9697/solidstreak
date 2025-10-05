package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/anatoliy9697/solidstreak/solidstreak-backend/pkg/date"
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
	Data *Habit    `json:"data"`
	Meta *Metadata `json:"meta"`
}

type PostPutHabitResponse struct {
	Data *hPkg.Habit `json:"data"`
}

type GetHabitsResponse struct {
	Data []*hPkg.Habit `json:"data"`
}

type HabitCheck struct {
	CheckDate *date.Date `json:"checkDate"`
	Completed *bool      `json:"completed"`
}

type PostHabitCheckRequest struct {
	Data *HabitCheck `json:"data"`
	Meta *Metadata   `json:"meta"`
}

type PostHabitCheckResponse struct {
	Data *hPkg.HabitCheck `json:"data"`
}

type GetUserHabitsCompletedChecksResponse struct {
	Data []*hPkg.HabitCheck `json:"data"`
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

	userTgID, ok := r.Context().Value(ctxKeyUserTgID{}).(int64)
	if !ok {
		err = apperrors.ErrUnauthorized("couldn't identify user")
		return
	}

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

	userTgID, ok := r.Context().Value(ctxKeyUserTgID{}).(int64)
	if !ok {
		err = apperrors.ErrUnauthorized("couldn't identify user")
		return
	}

	var req PostPutHabitRequest

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&req); err != nil {
		err = apperrors.ErrBadRequest("invalid request payload")
		return
	}

	if req.Data == nil {
		err = apperrors.ErrBadRequest("habit data is required")
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

	userTgID, ok := r.Context().Value(ctxKeyUserTgID{}).(int64)
	if !ok {
		err = apperrors.ErrUnauthorized("couldn't identify user")
		return
	}

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

	if req.Data == nil {
		err = apperrors.ErrBadRequest("habit data is required")
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

func (s Server) postUserHabitCheck(w http.ResponseWriter, r *http.Request) {
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

	var userID int64
	userID, err = getInt64FromURLParams(r, "userId", true)
	if err != nil {
		return
	}

	var habitID int64
	habitID, err = getInt64FromURLParams(r, "habitId", true)
	if err != nil {
		return
	}

	var req PostHabitCheckRequest

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&req); err != nil {
		err = apperrors.ErrBadRequest("invalid request payload")
		return
	}

	if req.Data == nil {
		err = apperrors.ErrBadRequest("request data is required")
		return
	}
	if req.Data.CheckDate == nil {
		err = apperrors.ErrBadRequest("habit check date is required")
		return
	}
	if req.Data.Completed == nil {
		err = apperrors.ErrBadRequest("habit completion status is required")
		return
	}

	var user *usrPkg.User
	if user, err = s.Res.UsrRepo.GetByTgID(userTgID); err != nil {
		return
	}

	if userID != user.ID {
		err = apperrors.ErrForbidden("couldn't check habit for another user")
		return
	}

	var habit *hPkg.Habit
	habit, err = s.Res.HabitRepo.GetByIDAndOwnerID(habitID, user.ID)
	if err != nil {
		return
	}

	habitCheck := &hPkg.HabitCheck{
		HabitID:   habit.ID,
		UserID:    user.ID,
		CheckDate: *req.Data.CheckDate,
		Completed: *req.Data.Completed,
		CheckedAt: time.Now(),
	}

	if err = s.Res.HabitRepo.SetUserHabitCheck(habitCheck); err != nil {
		return
	}

	response := PostHabitCheckResponse{Data: habitCheck}

	json.NewEncoder(w).Encode(response)
}

func (s Server) getUserHabitCompletedChecks(w http.ResponseWriter, r *http.Request) {
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

	var userID int64
	userID, err = getInt64FromURLParams(r, "userId", true)
	if err != nil {
		return
	}

	var habitID int64
	habitID, err = getInt64FromURLParams(r, "habitId", true)
	if err != nil {
		return
	}

	var fromDate, toDate *date.Date
	fromDate, err = getDateFromURLQuery(r, "from", false)
	if err != nil {
		return
	}
	if fromDate == nil {
		d := date.Today().AddDate(-1, 0, 0)
		fromDate = &d
	}
	toDate, err = getDateFromURLQuery(r, "to", false)
	if err != nil {
		return
	}
	if toDate == nil {
		d := date.Today()
		toDate = &d
	}

	var user *usrPkg.User
	if user, err = s.Res.UsrRepo.GetByTgID(userTgID); err != nil {
		return
	}

	// TODO: временно. будет исправлено в рамках https://github.com/anatoliy9697/solidstreak/issues/27
	if userID != user.ID {
		err = apperrors.ErrForbidden("couldn't get habit checks for another user")
		return
	}

	var habitChecks []*hPkg.HabitCheck
	if habitChecks, err = s.Res.HabitRepo.GetUserHabitsCompletedChecks(userID, []int64{habitID}, fromDate, toDate); err != nil {
		return
	}

	response := GetUserHabitsCompletedChecksResponse{Data: habitChecks}

	json.NewEncoder(w).Encode(response)
}
