package handler

import (
	"encoding/json"
	"fmt"
	"hack/internal/auth"
	"hack/internal/config"
	"hack/internal/models"
	"hack/internal/pkg"
	"hack/internal/pkg/errors"
	"io"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	pkgErrors "github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type handlers struct {
	authUC auth.UseCaseI
}

func New(aUC auth.UseCaseI) auth.HandlersI {
	return &handlers{
		authUC: aUC,
	}
}

func setNewCookie(w http.ResponseWriter, session *models.Session) {
	http.SetCookie(w, &http.Cookie{
		Name:     config.CookieName,
		Value:    session.SessionID,
		Expires:  time.Now().Add(config.CookieTTL),
		HttpOnly: true,
		Path:     config.CookiePath,
	})
}

func delCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    config.CookieName,
		Value:   "",
		Expires: time.Now().AddDate(0, 0, -1),
		Path:    config.CookiePath,
	})
}

func (h *handlers) SignUp(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error("failed to close request: ", err)
		}
	}(r.Body)

	form := models.FormSignUp{}
	if err := json.NewDecoder(r.Body).Decode(&form); err != nil {
		pkg.HandleError(w, r, pkgErrors.Wrap(errors.ErrInvalidForm, err.Error()))
		return
	}

	validate := validator.New()
	if err := validate.Struct(form); err != nil {
		pkg.HandleError(w, r, pkgErrors.Wrap(errors.ErrInvalidForm, err.Error()))
		return
	}

	response, session, err := h.authUC.SignUp(form)
	if err != nil {
		pkg.HandleError(w, r, err)
		return
	}

	setNewCookie(w, session)
	pkg.SendJSON(w, r, http.StatusOK, response)
}

func (h *handlers) SignIn(w http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Error(fmt.Errorf("failed to close request: %w", err))
		}
	}(r.Body)
	form := models.FormLogin{}

	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		pkg.HandleError(w, r, pkgErrors.Wrap(errors.ErrInvalidForm, err.Error()))
		return
	}

	response, session, err := h.authUC.SignIn(form)
	if err != nil {
		pkg.HandleError(w, r, err)
		return
	}

	setNewCookie(w, session)
	pkg.SendJSON(w, r, http.StatusOK, response)
}

func (h *handlers) Logout(w http.ResponseWriter, _ *http.Request) {
	delCookie(w)
	w.WriteHeader(http.StatusOK)
}
