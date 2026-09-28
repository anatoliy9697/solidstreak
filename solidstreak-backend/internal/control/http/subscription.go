package http

import (
	"encoding/json"
	"net/http"

	apperrors "github.com/anatoliy9697/solidstreak/solidstreak-backend/pkg/errors"

	subPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/subscription"
	usrPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/user"

	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/usecases"
)

type GetSubscriptionPlansResponse struct {
	Data []*subPkg.Plan `json:"data"`
}

type GetUserSubscriptionResponse struct {
	Data *subPkg.Subscription `json:"data"`
}

func (s Server) getSubscriptionPlans(w http.ResponseWriter, r *http.Request) {
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

	var plans []*subPkg.Plan
	if plans, err = s.Res.SubRepo.GetPlans(); err != nil {
		return
	}

	json.NewEncoder(w).Encode(GetSubscriptionPlansResponse{Data: plans})
}

func (s Server) getUserSubscription(w http.ResponseWriter, r *http.Request) {
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
		err = apperrors.NewUnauthorizedErr("couldn't identify user")
		return
	}

	var userID int64
	if userID, err = getInt64FromURLParams(r, "userID", true); err != nil {
		return
	}

	var user *usrPkg.User
	if user, err = s.Res.UsrRepo.GetByID(userID); err != nil {
		return
	}

	if user.TgID != userTgID {
		err = apperrors.NewUnauthorizedErr("couldn't get subscription for another user")
		return
	}

	if user.Subscription, err = usecases.GetUserSubscription(s.Res, user); err != nil {
		return
	}

	json.NewEncoder(w).Encode(GetUserSubscriptionResponse{Data: user.Subscription})
}
