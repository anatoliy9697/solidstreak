package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	apperrors "github.com/anatoliy9697/solidstreak/solidstreak-backend/pkg/errors"

	invPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/invoice"
	subPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/subscription"
	tcPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/tgchat"
	usrPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/user"
	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/usecases"
)

type Invoice struct {
	SubscriptionPlanCode    *string                        `json:"subscriptionPlanCode"`
	SubscriptionPeriodUnit  *subPkg.SubscriptionPeriodUnit `json:"subscriptionPeriodUnit"`
	SubscriptionPeriodCount *int64                         `json:"subscriptionPeriodCount"`
	Currency                *invPkg.Currency               `json:"currency"`
}

type PostInvoiceRequest struct {
	Data *Invoice `json:"data"`
}

type PostInvoiceResponse struct {
	Data *invPkg.Invoice `json:"data"`
}

func (s Server) postInvoice(w http.ResponseWriter, r *http.Request) {
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
	userID, err = getInt64FromURLParams(r, "userID", true)
	if err != nil {
		return
	}

	var req PostInvoiceRequest

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&req); err != nil {
		err = apperrors.NewBadRequestErr("invalid request payload")
		return
	}

	if req.Data == nil {
		err = apperrors.NewBadRequestErr("invoice data is required")
		return
	}
	if req.Data.SubscriptionPlanCode == nil {
		err = apperrors.NewBadRequestErr("subscription plan code is required")
		return
	}
	var plan *subPkg.Plan
	if plan, err = s.Res.SubRepo.GetPlanByCode(*req.Data.SubscriptionPlanCode); err != nil {
		return
	}
	if plan.Code == "basic" {
		err = apperrors.NewBadRequestErr("cannot purchase basic subscription")
		return
	}
	if req.Data.SubscriptionPeriodUnit == nil {
		err = apperrors.NewBadRequestErr("subscription period unit is required")
		return
	}
	var periodUnit subPkg.SubscriptionPeriodUnit
	if periodUnit, ok = subPkg.SubscriptionPeriodUnitMapping[string(*req.Data.SubscriptionPeriodUnit)]; !ok {
		err = apperrors.NewBadRequestErr("invalid subscription period unit")
		return
	}
	if req.Data.SubscriptionPeriodCount == nil {
		err = apperrors.NewBadRequestErr("subscription period count is required")
		return
	}
	if *req.Data.SubscriptionPeriodCount != 1 {
		err = apperrors.NewBadRequestErr("invalid subscription period count")
		return
	}
	if req.Data.Currency == nil {
		err = apperrors.NewBadRequestErr("currency is required")
		return
	}
	var currency invPkg.Currency
	if currency, ok = invPkg.CurrencyMapping[string(*req.Data.Currency)]; !ok {
		err = apperrors.NewBadRequestErr("invalid currency")
		return
	}

	var user *usrPkg.User
	if user, err = s.Res.UsrRepo.GetByID(userID); err != nil {
		return
	}

	if user.TgID != userTgID {
		err = apperrors.NewForbiddenErr("couldn't purchase subscription for another user")
		return
	}

	var tgChat *tcPkg.Chat
	if tgChat, err = s.Res.TCRepo.GetByUserID(userID); err != nil {
		return
	}

	var pricing *subPkg.Pricing
	for _, p := range plan.Pricing {
		if p.PeriodUnit == periodUnit && p.PeriodCount == *req.Data.SubscriptionPeriodCount && p.Currency == currency {
			pricing = &p
			break
		}
	}
	if pricing == nil {
		err = apperrors.NewNotFoundErr(fmt.Sprintf("couldn't find pricing for plan \"%s\", period unit \"%s\", period count \"%d\", currency \"%s\"", plan.Code, periodUnit, *req.Data.SubscriptionPeriodCount, currency))
		return
	}

	var subscription *subPkg.Subscription
	if subscription, err = usecases.GetUserSubscription(s.Res, user); err != nil {
		return
	}

	if subscription.PlanCode == "premium" && subscription.FinishDt == nil {
		err = apperrors.NewConflictErr("cannot purchase or renew a premium subscription while user already has an active lifetime subscription")
		return
	}

	var invoice *invPkg.Invoice
	if invoice, err = s.Res.InvRepo.GetActiveNotExpiredByStatusesAndUserID(userID, []string{string(invPkg.InvoiceStatusPending)}); err != nil && !errors.Is(err, apperrors.ErrNotFound) {
		return
	}
	if invoice != nil {
		err = apperrors.NewConflictErr("user already has an active and not expired invoice in pending status")
		return
	}

	if invoice, err = usecases.CreateAndSendInvoice(r.Context(), s.Res, user, tgChat, plan, periodUnit, *req.Data.SubscriptionPeriodCount, pricing); err != nil {
		return
	}

	response := PostInvoiceResponse{Data: invoice}

	json.NewEncoder(w).Encode(response)
}
