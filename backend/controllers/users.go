package controllers

import (
	"chatin/config/keys"
	"chatin/models"
	"chatin/server/request"
	"chatin/server/response"
	"chatin/services"
	"net/http"
)

type usersController struct{}

var usersService = services.NewUserService()

func NewUserController() *usersController {
	return &usersController{}
}

func (c *usersController) Get(res http.ResponseWriter, req *http.Request) {
	users, err := usersService.Get()
	if err != nil {
		return
	}
	response.JSON(res, users)
}

func (c *usersController) Create(res http.ResponseWriter, req *http.Request) {
	var user models.Users
	request.GetBody(req, &user)
	users, err := usersService.Create(user)
	if err != nil {
		return
	}
	response.JSON(res, users)
}

func (c *usersController) Find(res http.ResponseWriter, req *http.Request) {
	var id = req.Context().Value(keys.UserContextKey).(float64)
	if id == 0 {
		http.Error(res, "unanthenticated", http.StatusUnauthorized)
		return
	}
	var user, err = usersService.Find(int(id))
	if err != nil {
		http.Error(res, "USER NOT FOUND", 404)
	}
	response.JSON(res, user)

}
