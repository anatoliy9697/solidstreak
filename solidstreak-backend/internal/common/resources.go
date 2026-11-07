package common

import (
	"log/slog"

	tgbotapi "github.com/mymmrac/telego"

	h "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/habit/repo"
	inv "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/invoice/repo"
	sub "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/subscription/repo"
	st "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/supportticket/repo"
	tr "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/terms/repo"
	tc "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/tgchat/repo"
	usr "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/user/repo"
)

type Resources struct {
	WebAppURL string
	Logger    *slog.Logger
	TgBotAPI  *tgbotapi.Bot
	UsrRepo   usr.Repo
	SubRepo   sub.Repo
	TCRepo    tc.Repo
	HabitRepo h.Repo
	InvRepo   inv.Repo
	STRepo    st.Repo
	TermsRepo tr.Repo
}
