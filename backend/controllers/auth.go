package controllers

import (
	"chatin/models"
	"chatin/server/request"
	"chatin/server/response"
	"chatin/services"
	"net/http"
)

type authControlller struct{}

func NewAuthController() *authControlller {
	return &authControlller{}
}

var authService = services.NewAuthService()

func (c *authControlller) Login(res http.ResponseWriter, req *http.Request) {
	var userData models.Users
	request.GetBody(req, &userData)
	var token, err = authService.Login(userData.Email, userData.Password)
	if err != nil {
		res.Header().Set("Content-Type", "application/json")
		http.Error(res, err.Error(), 500)
	}
	response.JSON(res, token)
}
