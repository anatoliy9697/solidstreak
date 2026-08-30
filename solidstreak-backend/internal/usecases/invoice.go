package usecases

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"context"

	"github.com/google/uuid"
	tgbotapi "github.com/mymmrac/telego"

	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/common"
	invPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/invoice"
	subPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/subscription"
	tcPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/tgchat"
	usrPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/user"
)

func CreateAndSendInvoice(
	ctx context.Context,
	r common.Resources,
	u *usrPkg.User,
	tc *tcPkg.Chat,
	plan *subPkg.Plan,
	subscriptionPeriodUnit subPkg.SubscriptionPeriodUnit,
	subscriptionPeriodCount int64,
	pricing *subPkg.Pricing,
	invoiceExpiresIn time.Duration,
) (*invPkg.Invoice, error) {
	var err error

	invoice := invPkg.NewInvoice(
		uuid.NewString(), pricing.Currency, pricing.Price, u.ID,
		time.Now().Add(invoiceExpiresIn),
	)

	subscriptionEvent := subPkg.NewSubscriptionEvent(
		subPkg.SubscriptionEventTypeAcquisition,
		subPkg.SubscriptionEventStatusInProgress,
		subPkg.SubscriptionOriginPurchase,
		plan.Code,
		subscriptionPeriodUnit,
		subscriptionPeriodCount,
		u.ID,
		invoice.UUID,
	)

	var tgInvoiceParams *tgbotapi.SendInvoiceParams
	if tgInvoiceParams, err = NewTgInvoiceParams(r, tc.TgID, u.LangCode, invoice, subscriptionEvent, invoiceExpiresIn); err != nil {
		return nil, err
	}

	var tgInvoiceMsg *tgbotapi.Message
	if tgInvoiceMsg, err = SendTgInvoice(ctx, r.TgBotAPI, tgInvoiceParams); err != nil {
		return nil, err
	}

	invoice.SetTgMessageID(tgInvoiceMsg.MessageID)

	if err = r.InvRepo.Create(invoice); err != nil {
		return nil, err
	}

	if err = r.SubRepo.CreateEvent(subscriptionEvent); err != nil {
		return nil, err
	}

	return invoice, nil
}

func BuildInvoicePayloadString(invoice *invPkg.Invoice, subscriptionEvent *subPkg.SubscriptionEvent) string {
	return "user:" + fmt.Sprint(invoice.UserID) + ":subscription:" + subscriptionEvent.SubscriptionPlanCode + ":periodUnit:" + string(subscriptionEvent.SubscriptionPeriodUnit) + ":periodCount:" + fmt.Sprint(subscriptionEvent.SubscriptionPeriodCount) + ":currency:" + string(invoice.Currency) + ":invoiceUUID:" + invoice.UUID
}

func ExtractInvoiceUUIDFromPayload(invoicePayload string) (string, error) {
	payloadParts := strings.Split(invoicePayload, ":")
	if len(payloadParts) != 12 {
		return "", errors.New("invalid invoice payload")
	}

	return payloadParts[11], nil
}

func MarkInvoiceAsPaid(r common.Resources, invoice *invPkg.Invoice, telegramPaymentChargeID string) error {
	invoice.MarkAsPaid(telegramPaymentChargeID)
	return r.InvRepo.Update(invoice)
}
