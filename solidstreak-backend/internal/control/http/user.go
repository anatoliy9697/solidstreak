package http

import (
	"encoding/json"
	"net/http"
	"net/url"

	apperrors "github.com/anatoliy9697/solidstreak/solidstreak-backend/pkg/errors"

	tcPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/tgchat"
	usrPkg "github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/domain/user"
)

type InitDataUser struct {
	TgID        int64  `json:"id"`
	TgUsername  string `json:"username"`
	TgFirstName string `json:"first_name"`
	TgLastName  string `json:"last_name"`
	TgLangCode  string `json:"language_code"`
	TgIsBot     bool   `json:"is_bot"`
}

type InitDataTgChat struct {
	TgID int64 `json:"id"`
}

type InputUser struct {
	TgID        int64  `json:"tgId"`
	TgUsername  string `json:"tgUsername"`
	TgFirstName string `json:"tgFirstName"`
	TgLastName  string `json:"tgLastName"`
	TgLangCode  string `json:"tgLangCode"`
	TgIsBot     bool   `json:"tgIsBot"`
}

type TgChat struct {
	TgID int64 `json:"tgId"`
}

type UserInfoData struct {
	User   *InputUser `json:"user"`
	TgChat *TgChat    `json:"tgChat"`
}

type PostUserInfoRequest struct {
	Data *UserInfoData `json:"data"`
}

type GetUserResponse struct {
	Data *usrPkg.User `json:"data"`
}

func (s Server) postUserInfo(w http.ResponseWriter, r *http.Request) {
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

	var req PostUserInfoRequest

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&req); err != nil {
		err = apperrors.ErrBadRequest("invalid request payload")
		return
	}

	if req.Data == nil {
		err = apperrors.ErrBadRequest("request data is required")
		return
	}

	inputUser := req.Data.User
	if inputUser == nil {
		err = apperrors.ErrBadRequest("user data is required")
		return
	}

	inputTgChat := req.Data.TgChat
	if inputTgChat == nil {
		err = apperrors.ErrBadRequest("chat data is required")
		return
	}

	var (
		initDataUser   *InitDataUser
		initDataTgChat *InitDataTgChat
	)
	if initDataUser, initDataTgChat, err = getUserAndChatFromInitData(r.Header.Get("X-Telegram-InitData")); err != nil {
		return
	}

	if initDataTgChat == nil {
		initDataTgChat = &InitDataTgChat{TgID: initDataUser.TgID} // Use personal chat with user if no other chat info
	}

	if initDataUser.TgID != inputUser.TgID || initDataUser.TgUsername != inputUser.TgUsername ||
		initDataUser.TgFirstName != inputUser.TgFirstName || initDataUser.TgLastName != inputUser.TgLastName ||
		initDataUser.TgLangCode != inputUser.TgLangCode || initDataUser.TgIsBot != inputUser.TgIsBot {
		err = apperrors.ErrBadRequest("user data does not match init data")
		return
	}

	if initDataTgChat.TgID != inputTgChat.TgID {
		err = apperrors.ErrBadRequest("telegram chat data does not match init data")
		return
	}

	user := usrPkg.NewUser(
		inputUser.TgID,
		inputUser.TgUsername,
		inputUser.TgFirstName,
		inputUser.TgLastName,
		inputUser.TgLangCode,
		inputUser.TgIsBot,
	)

	userExists := false
	if userExists, err = s.Res.UsrRepo.IsExists(user); err != nil {
		return
	}
	if userExists {
		err = s.Res.UsrRepo.Update(user)
	} else {
		err = s.Res.UsrRepo.Create(user)
	}
	if err != nil {
		return
	}

	tgChat := tcPkg.NewChat(inputTgChat.TgID, user.ID)

	chatExists := false
	if chatExists, err = s.Res.TCRepo.IsExistsByTgID(tgChat.TgID); err != nil {
		return
	}

	if chatExists {
		err = s.Res.TCRepo.Update(tgChat)
	} else {
		err = s.Res.TCRepo.Create(tgChat)
	}
	if err != nil {
		return
	}

	json.NewEncoder(w).Encode(GetUserResponse{Data: user})
}

func (s Server) getUser(w http.ResponseWriter, r *http.Request) {
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

	var userId int64
	if userId, err = getInt64FromURLParams(r, "userId", true); err != nil {
		return
	}

	var user *usrPkg.User
	if user, err = s.Res.UsrRepo.GetByID(userId); err != nil {
		return
	}

	// Временное ограничение. Как появится возможность просмтатривать карточки чужих юзеров, убрать это условие
	if user.TgID != userTgID {
		err = apperrors.ErrUnauthorized("couldn't get user info for another user")
		return
	}

	json.NewEncoder(w).Encode(GetUserResponse{Data: user})
}

func getUserAndChatFromInitData(initData string) (*InitDataUser, *InitDataTgChat, error) {
	var (
		user *InitDataUser
		chat *InitDataTgChat
	)

	values, err := url.ParseQuery(initData)
	if err != nil {
		return nil, nil, err
	}

	for k, v := range values {
		switch k {
		case "user":
			if err := json.Unmarshal([]byte(v[0]), &user); err != nil {
				return nil, nil, err
			}
		case "chat":
			if err := json.Unmarshal([]byte(v[0]), &chat); err != nil {
				return nil, nil, err
			}
		}
		if user != nil && chat != nil {
			break
		}
	}

	return user, chat, nil
}
