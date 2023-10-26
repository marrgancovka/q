package handler

import (
	"hack/internal/auth"
	"hack/internal/config"
	middleware "hack/internal/middlware"
	"net/http"

	"github.com/gorilla/mux"
)

func RegisterHTTPRoutes(r *mux.Router, h auth.HandlersI, m *middleware.Middleware) {
	r.HandleFunc(config.RouteSignin, h.SignIn).Methods(http.MethodPost)
	r.HandleFunc(config.RouteSignup, h.SignUp).Methods(http.MethodPost)
	r.HandleFunc(config.RouteLogout, m.CheckAuth(h.Logout)).Methods(http.MethodGet)
}
