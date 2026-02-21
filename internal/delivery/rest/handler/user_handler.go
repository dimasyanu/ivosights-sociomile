package handler

import "github.com/dimasyanu/ivosights-sociomile/service"

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}
