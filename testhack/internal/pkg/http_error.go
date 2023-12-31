package pkg

import (
	"net/http"

	"hack/internal/pkg/errors"

	pkgErr "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	logger, ok := r.Context().Value(ContextHandlerLog).(*Logger)
	if !ok {
		log.Error("failed to get logger for handler", r.URL.Path)
		log.Error(err)
	} else {
		logger.Error(err)
	}

	causeErr := pkgErr.Cause(err)
	code := errors.Code(causeErr)
	customErr := errors.New(code, causeErr)
	SendJSON(w, r, code, customErr)
}
