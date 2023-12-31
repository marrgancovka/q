package usecase

import (
	"fmt"
	"hack/internal/auth"
	"hack/internal/models"
	"hack/internal/user"

	pkgErrors "github.com/pkg/errors"
)

type sessionUC struct {
	sessionRepo auth.SessionRepoI
	userUC      user.UseCaseI
}

func NewSessionUC(sr auth.SessionRepoI, uc user.UseCaseI) auth.SessionUseCaseI {
	return &sessionUC{
		sessionRepo: sr,
		userUC:      uc,
	}
}

func (u *sessionUC) CreateSession(uID uint64) (*models.Session, error) {
	newSession, err := u.sessionRepo.CreateSession(uID)
	if err != nil {
		return nil, pkgErrors.Wrap(err, "create session")
	}

	return newSession, nil
}

func (u *sessionUC) DeleteSession(sessionID string) error {
	err := u.sessionRepo.DeleteSession(sessionID)
	if err != nil {
		return pkgErrors.Wrap(err, "delete avatar")
	}

	return nil
}

func (u *sessionUC) GetSession(sessionID string) (*models.Session, error) {
	s, err := u.sessionRepo.GetSession(sessionID)
	if err != nil {
		return nil, fmt.Errorf("get session: %w", err)
	}

	return s, nil
}
