package models

type AuthResponse struct {
	Email     string `json:"email" validate:"required"`
	FirstName string `json:"firstName" validate:"required"`
	LastName  string `json:"lastName" validate:"required"`
}

type FormSignUp struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
	RepeatPw string `json:"repeatPw" validate:"required"`
	// FirstName string `json:"firstName" validate:"required"`
	// LastName  string `json:"lastName" validate:"required"`
}

type FormLogin struct {
	Login    string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required"`
	// Remember bool   `json:"remember" validate:"required"`
}
