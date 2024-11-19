package handler

import (
	"net/http"

	"server/internal/entity"
	"server/internal/helper"
	"server/internal/service"

	"github.com/gin-gonic/gin"
)

type usersHandler struct {
	usersService service.UsersService
}

func NewUsersHandler(usersService service.UsersService) *usersHandler {
	return &usersHandler{usersService}
}

func (user_handler *usersHandler) Register(c *gin.Context) {
	var userRequest entity.UsersRequest
	err := c.ShouldBindBodyWithJSON(&userRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	user, err := user_handler.usersService.Register(userRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	// MENGUBAH TIPE ENITITY KE TIPE RESPONSE
	userResponse := helper.ConvertToResponseType(user)

	c.JSON(http.StatusCreated, gin.H{
		"statusCode": 201,
		"status":     true,
		"message":    "Register user data",
		"data":       userResponse,
	})
}

func (user_handler *usersHandler) Login(c *gin.Context) {
	var userRequest entity.UsersRequest
	err := c.ShouldBindBodyWithJSON(&userRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	token, err := user_handler.usersService.Login(userRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	c.SetCookie("token", *token, 60*60*24, "/", "localhost", false, false)

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"status":     true,
		"message":    "Login user data",
		"data":       token,
	})
}

func (user_handler *usersHandler) GetAllUsers(c *gin.Context) {
	users, err := user_handler.usersService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	// MENGUBAH TIPE ENITITY KE TIPE RESPONSE
	var usersResponse []entity.UsersResponse
	for _, user := range users {
		userResponse, _ := helper.ConvertToResponseType(user).(entity.UsersResponse)
		usersResponse = append(usersResponse, userResponse)
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"status":     true,
		"message":    "Get all users data",
		"data":       usersResponse,
	})
}

func (user_handler *usersHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	user, err := user_handler.usersService.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	// MENGUBAH TIPE ENITITY KE TIPE RESPONSE
	userResponse := helper.ConvertToResponseType(user)

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"status":     true,
		"message":    "Get user data",
		"data":       userResponse,
	})
}

func (user_handler *usersHandler) UpdateUser(c *gin.Context) {
	var userRequest entity.UsersRequest
	err := c.ShouldBindBodyWithJSON(&userRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	id := c.Param("id")

	user, err := user_handler.usersService.UpdateUser(id, userRequest)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	// MENGUBAH TIPE ENITITY KE TIPE RESPONSE
	userResponse := helper.ConvertToResponseType(user)

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"status":     true,
		"message":    "Update user data",
		"data":       userResponse,
	})
}

func (user_handler *usersHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")

	user, err := user_handler.usersService.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	// MENGUBAH TIPE ENITITY KE TIPE RESPONSE
	userResponse := helper.ConvertToResponseType(user)

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"status":     true,
		"message":    "Delete user data",
		"data":       userResponse,
	})
}
