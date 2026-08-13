package http

import (
	"encoding/json"
	"net/http"

	subPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/subscription"
)

type GetSubscriptionPlansResponse struct {
	Data []*subPkg.Plan `json:"data"`
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
