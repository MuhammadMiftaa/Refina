package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"golang.org/x/oauth2"

	"server/internal/entity"
	"server/internal/helper"
	"server/internal/service"

	"github.com/gin-gonic/gin"
)

var mode = os.Getenv("MODE")

type usersHandler struct {
	usersService service.UsersService
	otpService   service.OTPService
}

func NewUsersHandler(usersService service.UsersService, otpService service.OTPService) *usersHandler {
	return &usersHandler{
		usersService: usersService,
		otpService:   otpService,
	}
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

func (user_handler *usersHandler) OAuthHandler(state string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			config *oauth2.Config
			err    error
		)
		switch state {
		case "google":
			config, _, err = helper.GetGoogleOAuthConfig()
		case "github":
			config, _, err = helper.GetGithubOAuthConfig()
		case "microsoft":
			config, _, err = helper.GetMicrosoftOAuthConfig()
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"statusCode": 500,
				"status":     false,
				"message":    err.Error(),
			})
			return
		}

		url := config.AuthCodeURL(state, oauth2.AccessTypeOffline) // BESERTA REFRESH TOKEN
		// c.Redirect(http.StatusFound, url) // VIA BACKEND
		c.JSON(http.StatusOK, gin.H{"url": url}) // VIA FRONTEND
	}
}

func (user_handler *usersHandler) CallbackGoogle(c *gin.Context) {
	// Ambil konfigurasi OAuth Google
	config, redirect_url, err := helper.GetGoogleOAuthConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": 500,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	// Ambil authorization code dari query parameter
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code not found"})
		return
	}

	// Tukar authorization code dengan access token
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	// Gunakan access token untuk mengambil informasi pengguna
	client := config.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	// Parse data pengguna
	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user info"})
		return
	}

	tokenJWT, err := user_handler.usersService.OAuthLogin(userInfo["name"].(string), userInfo["email"].(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	if mode == "production" {
		c.Redirect(http.StatusFound, redirect_url+"/login?token="+*tokenJWT)
	}
	c.SetCookie("token", *tokenJWT, 60*60*24, "/", "localhost", false, false)

	c.Redirect(http.StatusFound, redirect_url)
}

func (user_handler *usersHandler) CallbackGithub(c *gin.Context) {
	// Ambil konfigurasi OAuth Google
	config, redirect_url, err := helper.GetGithubOAuthConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": 500,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	// Ambil authorization code dari query parameter
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code not found"})
		return
	}

	// Tukar authorization code dengan access token
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	// Gunakan access token untuk mengambil informasi pengguna
	client := config.Client(context.Background(), token)
	// Ambil data pengguna
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	// Ambil email pengguna
	emailResp, err := client.Get("https://api.github.com/user/emails")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user email"})
		return
	}
	defer emailResp.Body.Close()

	// Baca data dari io.ReadCloser (resp.Body)
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read user info"})
		return
	}

	var githubUser helper.GitHubUser
	if err := json.Unmarshal(data, &githubUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user info"})
		return
	}

	// Parse email data
	var emails []map[string]interface{}
	if err := json.NewDecoder(emailResp.Body).Decode(&emails); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse email data"})
		return
	}

	// Pilih email utama (primary)
	var primaryEmail string
	for _, email := range emails {
		if isPrimary, ok := email["primary"].(bool); ok && isPrimary {
			if emailAddress, ok := email["email"].(string); ok {
				primaryEmail = emailAddress
				break
			}
		}
	}

	tokenJWT, err := user_handler.usersService.OAuthLogin(githubUser.Name, primaryEmail)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	if mode == "production" {
		c.Redirect(http.StatusFound, redirect_url+"/login?token="+*tokenJWT)
	}
	c.SetCookie("token", *tokenJWT, 60*60*24, "/", "localhost", false, false)

	c.Redirect(http.StatusFound, redirect_url)
}

func (user_handler *usersHandler) CallbackMicrosoft(c *gin.Context) {
	// Ambil konfigurasi OAuth Google
	config, redirect_url, err := helper.GetMicrosoftOAuthConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": 500,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	// Ambil authorization code dari query parameter
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Authorization code not found"})
		return
	}

	// Tukar authorization code dengan access token
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	// Gunakan access token untuk mengambil informasi pengguna
	client := config.Client(context.Background(), token)
	// Ambil data pengguna
	resp, err := client.Get("https://graph.microsoft.com/v1.0/me")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	// Parse data pengguna
	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse user info"})
		return
	}

	tokenJWT, err := user_handler.usersService.OAuthLogin(userInfo["displayName"].(string), userInfo["mail"].(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	if mode == "production" {
		c.Redirect(http.StatusFound, redirect_url+"/login?token="+*tokenJWT)
	}
	c.SetCookie("token", *tokenJWT, 60*60*24, "/", "localhost", false, false)

	c.Redirect(http.StatusFound, redirect_url)
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

func (user_handler *usersHandler) SendOTP(c *gin.Context) {
	var OTP helper.OTP
	if err := c.ShouldBindBodyWithJSON(&OTP); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"status":     false,
			"message":    "Invalid request body",
		})
		return
	}
	OTP.OTP = helper.GenerateOTP()

	// Simpan OTP ke Redis
	if err := user_handler.otpService.SetOTP(OTP.Email, OTP.OTP, 5*time.Minute); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": 500,
			"status":     false,
			"message":    "Failed to save OTP",
		})
		return
	}

	// Kirimkan OTP ke email
	if err := helper.SendEmail(OTP.Email, OTP.OTP); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to send email"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": 200,
		"message":    "OTP sent successfully",
		"data":       OTP.Email,
	})
}

func (user_handler *usersHandler) VerifyOTP(c *gin.Context) {
	var OTP helper.OTP
	if err := c.ShouldBindBodyWithJSON(&OTP); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	valid, err := user_handler.otpService.ValidateOTP(OTP.Email, OTP.OTP)
	if err != nil || !valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"statusCode": 401,
			"status":     false,
			"message":    "Invalid or expired OTP",
		})
		return
	}

	user, err := user_handler.usersService.VerifyUser(OTP.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"statusCode": 500,
			"status":     false,
			"message":    "Failed to verify user",
		})
		return
	}

	userResponse := helper.ConvertToResponseType(user)

	c.JSON(http.StatusOK, gin.H{
		"status":     true,
		"statusCode": 200,
		"message":    "OTP verified successfully",
		"data":       userResponse,
	})
}

func (user_handler *usersHandler) GetUserWallets(c *gin.Context) {
	id := c.Param("id")

	userWallets, err := user_handler.usersService.GetUserWallets(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"status":     true,
		"message":    "Get user wallets data",
		"data":       userWallets,
	})
}

func (user_handler *usersHandler) GetUserInvestments(c *gin.Context) {
	id := c.Param("id")

	userInvestments, err := user_handler.usersService.GetUserInvestments(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"status":     true,
		"message":    "Get user investments data",
		"data":       userInvestments,
	})
}

func (user_handler *usersHandler) GetUserTransactions(c *gin.Context) {
	id := c.Param("id")

	userTransactions, err := user_handler.usersService.GetUserTransactions(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"statusCode": 400,
			"status":     false,
			"message":    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"statusCode": 200,
		"status":     true,
		"message":    "Get user transactions data",
		"data":       userTransactions,
	})
}