package http

import (
	"net/http"

	"github.com/anatoliy9697/solidstreak/solidstreak-backend/internal/common/i18n"
)

func (s Server) getTerms(w http.ResponseWriter, r *http.Request) {
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

	var lang string
	lang, err = getStringFromURLParams(r, "lang", false)
	if err != nil {
		return
	}

	var version string
	version, err = getStringFromURLParams(r, "version", false)
	if err != nil {
		return
	}

	if lang == "" {
		lang = i18n.GetDefaultLang()
	}

	var doc []byte
	if doc, err = s.Res.TermsRepo.GetByLangAndVersion(lang, version); err != nil {
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	w.Write(doc)
}
