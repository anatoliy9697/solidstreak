package scheduler

import (
	"context"
	"fmt"

	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/common"

	invPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/invoice"
	st "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/schedulertask"
	subPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/subscription"
	usrPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/user"

	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/usecases"
)

type TaskRunner struct {
	ID  string
	Res common.Resources
}

func (tr TaskRunner) Run(ctx context.Context, doneCh chan string, tasks []st.Task) {
	defer func() { doneCh <- tr.ID }()

	tr.Res.Logger = tr.Res.Logger.With("taskRunnerID", tr.ID)

	tr.Res.Logger.Debug("task runner started", "taskCount", len(tasks))

	var task st.Task
	for _, task = range tasks {
		switch current := task.(type) {
		case invPkg.ProcessExpiredInvoiceTask:
			tr.processExpiredInvoice(ctx, current)
		default:
			tr.Res.Logger.Error(fmt.Sprintf("unknown task type: %T", current))
		}
	}

	tr.Res.Logger.Debug("task runner execution completed")
}

func (tr TaskRunner) processExpiredInvoice(ctx context.Context, task invPkg.ProcessExpiredInvoiceTask) error {
	var err error

	defer func() {
		if r := recover(); r != nil {
			tr.Res.Logger.Error("panic recovered in TaskRunner.processExpiredInvoice", "panic", r)
		}

		if err != nil {
			tr.Res.Logger.Error("task runner error", "error", err)
		}
	}()

	i := task.Invoice

	tr.Res.Logger = tr.Res.Logger.With("invoiceUUID", i.UUID)
	tr.Res.Logger.Debug("processing expired invoice")

	var se *subPkg.SubscriptionEvent
	if se, err = tr.Res.SubRepo.GetActiveEventByInvoiceUUIDAndStatuses(i.UUID, []subPkg.SubscriptionEventStatus{subPkg.SubscriptionEventStatusInProgress}); err != nil {
		return err
	}

	i.SetStatus(invPkg.InvoiceStatusExpired)

	if err = tr.Res.InvRepo.Update(i); err != nil {
		return err
	}

	se.SetStatus(subPkg.SubscriptionEventStatusPaymentTimeout)

	if err = tr.Res.SubRepo.UpdateEvent(se); err != nil {
		return err
	}

	var u *usrPkg.User
	if u, err = tr.Res.UsrRepo.GetByID(se.UserID); err != nil {
		return err
	}

	deleteMsgParams := usecases.NewDeleteTgMessageParams(i.TgChatID, i.TgMessageID)
	if err = usecases.DeleteTgMessage(ctx, tr.Res.TgBotAPI, deleteMsgParams); err != nil {
		return err
	}

	var lang string
	if u != nil {
		lang = u.LangCode
	} else {
		lang = common.GetDefaultLang()
	}

	if err = usecases.SendTgMessage(ctx, tr.Res.TgBotAPI, usecases.NewExpiredInvoiceTgMessageParams(tr.Res, i.TgChatID, lang)); err != nil {
		return err
	}

	tr.Res.Logger.Debug("expired invoice processed successfully")

	return nil
}
