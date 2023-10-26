package auth

import (
	"hack/internal/models"
	"net/http"
)

type HandlersI interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	SignIn(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
}

type UseCaseI interface {
	SignIn(form models.FormLogin) (*models.AuthResponse, *models.Session, error)
	SignUp(form models.FormSignUp) (*models.AuthResponse, *models.Session, error)
}

type SessionUseCaseI interface {
	CreateSession(uID uint64) (*models.Session, error)
	DeleteSession(sessionID string) error
	GetSession(sessionID string) (*models.Session, error)
}

type SessionRepoI interface {
	CreateSession(uID uint64) (*models.Session, error)
	DeleteSession(sessionID string) error
	GetSession(sessionID string) (*models.Session, error)
}
