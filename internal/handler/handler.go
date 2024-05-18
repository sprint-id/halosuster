package handler

import (
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/sprint-id/eniqilo-server/internal/cfg"
	"github.com/sprint-id/eniqilo-server/internal/service"
)

type Handler struct {
	router  *chi.Mux
	service *service.Service
	cfg     *cfg.Cfg
}

func NewHandler(router *chi.Mux, service *service.Service, cfg *cfg.Cfg) *Handler {
	handler := &Handler{router, service, cfg}
	handler.registRoute()

	return handler
}

func (h *Handler) registRoute() {

	r := h.router
	var tokenAuth *jwtauth.JWTAuth = jwtauth.New("HS256", []byte(h.cfg.JWTSecret), nil, jwt.WithAcceptableSkew(30*time.Second))

	userH := newUserHandler(h.service.User)
	patientH := newPatientHandler(h.service.Patient)
	recordH := newRecordHandler(h.service.Record)
	fileH := newFileHandler(h.cfg)

	r.Use(middleware.RedirectSlashes)

	r.Post("/v1/user/it/register", userH.RegisterIT)
	r.Post("/v1/user/it/login", userH.LoginIT)
	r.Post("/v1/user/nurse/login", userH.LoginNurse)

	// protected route
	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator(tokenAuth))

		r.Post("/v1/user/nurse/register", userH.RegisterNurse)
		r.Get("/v1/user", userH.GetUser)
		r.Put("/v1/user/nurse/{userId}", userH.UpdateNurse)
		r.Delete("/v1/user/nurse/{userId}", userH.DeleteNurse)
		r.Post("/v1/user/nurse/{userId}/access", userH.AccessNurse)

		r.Post("/v1/medical/patient", patientH.CreatePatient)
		r.Get("/v1/medical/patient", patientH.GetPatient)
		r.Post("/v1/medical/record", recordH.AddRecord)
		r.Get("/v1/medical/record", recordH.GetRecord)

		r.Post("/v1/image", fileH.Upload)
	})
}
