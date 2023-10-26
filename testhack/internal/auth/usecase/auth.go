package usecase

import (
	"hack/internal/auth"
	"hack/internal/config"
	"hack/internal/models"
	"hack/internal/pkg"
	"hack/internal/pkg/errors"
	"hack/internal/user"

	pkgErrors "github.com/pkg/errors"
)

type authUC struct {
	sessionUC auth.SessionUseCaseI
	userUC    user.UseCaseI
}

func NewAuthUC(sUC auth.SessionUseCaseI, uc user.UseCaseI) auth.UseCaseI {
	return &authUC{
		sessionUC: sUC,
		userUC:    uc,
	}
}

func (u *authUC) SignIn(form models.FormLogin) (*models.AuthResponse, *models.Session, error) {
	email, err := pkg.ValidateLogin(form.Login)
	if err != nil {
		return nil, nil, errors.ErrInvalidLogin
	}
	user, err := u.userUC.GetByEmail(email)
	if err != nil {
		return nil, nil, errors.ErrWrongPw
	}

	if user.Password != form.Password {
		return nil, nil, errors.ErrWrongPw
	}

	session, err := u.sessionUC.CreateSession(user.UserID)
	if err != nil {
		return nil, nil, pkgErrors.Wrap(err, "sign in")
	}

	return &models.AuthResponse{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, session, nil
}

func (u *authUC) SignUp(form models.FormSignUp) (*models.AuthResponse, *models.Session, error) {
	if form.RepeatPw != form.Password {
		return nil, nil, errors.ErrPwDontMatch
	}

	email, err := pkg.ValidateLogin(form.Login)
	if err != nil || len(form.Login) > 30 || len(form.Login) < 3 {
		return nil, nil, errors.ErrInvalidLogin
	}

	user, err := u.userUC.Create(&models.User{
		Email:    email,
		Password: form.Password,
		Avatar:   config.DefaultAvatar,
	})
	if err != nil {
		return nil, nil, pkgErrors.Wrap(err, "sign up")
	}

	session, err := u.sessionUC.CreateSession(user.UserID)
	if err != nil {
		//надо ли тут откатить прошлое ?
		return nil, nil, pkgErrors.Wrap(err, "sign up")
	}

	return &models.AuthResponse{
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}, session, nil
}
