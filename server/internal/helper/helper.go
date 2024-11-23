package helper

import (
	"errors"
	"fmt"
	"math/rand"
	"net/smtp"
	"os"
	"regexp"
	"time"
	"unicode"

	"server/internal/entity"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/microsoft"
)

func EmailValidator(str string) bool {
	email_validator := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return email_validator.MatchString(str)
}

func PasswordValidator(str string) (bool, bool, bool) {
	var hasLetter, hasDigit, hasMinLen bool
	for _, char := range str {
		switch {
		case unicode.IsLetter(char):
			hasLetter = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}

	if len(str) >= 8 {
		hasMinLen = true
	}

	return hasMinLen, hasLetter, hasDigit
}

func PasswordHashing(str string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(hashPassword), nil
}

func ConvertToResponseType(data interface{}) interface{} {
	switch v := data.(type) {
	case entity.Users:
		return entity.UsersResponse{
			ID:    v.ID.String(),
			Name:  v.Name,
			Email: v.Email,
		}
	case entity.Transactions:
		return entity.TransactionsResponse{
			ID:              v.ID.String(),
			Amount:          v.Amount,
			TransactionType: v.TransactionType,
			Date:            v.Date,
			Description:     v.Description,
			Category:        v.Category,
			UserID:          v.UserID,
		}
	default:
		return nil
	}
}

var secretKey = "pojq09720ef1ko0f1h9iego2010j20240"

func GenerateToken(username string, email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"username": username,
		"email":    email,
		"exp":      expirationTime.Unix(),
	}

	parseToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := parseToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func VerifyToken(cookie string) (interface{}, error) {
	token, _ := jwt.Parse(cookie, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("sign in to preceed")
		}
		return []byte(secretKey), nil
	})

	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		return nil, errors.New("sign in to preceed")
	}

	return token.Claims.(jwt.MapClaims), nil
}

func ComparePass(hashPassword, reqPassword string) bool {
	hash, pass := []byte(hashPassword), []byte(reqPassword)

	err := bcrypt.CompareHashAndPassword(hash, pass)
	return err == nil
}

func StorageIsExist(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

func GetGoogleOAuthConfig() (*oauth2.Config, string, error) {
	if err := godotenv.Load(); err != nil {
		return nil, "", err
	}

	var (
		ClientID          = os.Getenv("GOOGLE_CLIENT_ID")
		ClientSecret      = os.Getenv("GOOGLE_CLIENT_SECRET")
		redirectURL       = os.Getenv("FRONTEND_URL")
		googleOauthConfig = &oauth2.Config{
			ClientID:     ClientID,
			ClientSecret: ClientSecret,
			RedirectURL:  "http://localhost:8080/v1/auth/callback/google",
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		}
	)

	return googleOauthConfig, redirectURL, nil
}

func GetGithubOAuthConfig() (*oauth2.Config, string, error) {
	if err := godotenv.Load(); err != nil {
		return nil, "", err
	}

	var (
		ClientID          = os.Getenv("GITHUB_CLIENT_ID")
		ClientSecret      = os.Getenv("GITHUB_CLIENT_SECRET")
		redirectURL       = os.Getenv("FRONTEND_URL")
		githubOauthConfig = &oauth2.Config{
			ClientID:     ClientID,
			ClientSecret: ClientSecret,
			RedirectURL:  "http://localhost:8080/v1/auth/callback/github",
			Scopes: []string{
				"read:user",
				"user:email",
			},
			Endpoint: github.Endpoint,
		}
	)

	return githubOauthConfig, redirectURL, nil
}

func GetMicrosoftOAuthConfig() (*oauth2.Config, string, error) {
	if err := godotenv.Load(); err != nil {
		return nil, "", err
	}

	var (
		ClientID             = os.Getenv("MICROSOFT_CLIENT_ID")
		ClientSecret         = os.Getenv("MICROSOFT_CLIENT_SECRET")
		redirectURL          = os.Getenv("FRONTEND_URL")
		microsoftOauthConfig = &oauth2.Config{
			ClientID:     ClientID,
			ClientSecret: ClientSecret,
			RedirectURL:  "http://localhost:8080/v1/auth/callback/microsoft",
			Scopes: []string{
				"User.Read",
			},
			Endpoint: microsoft.AzureADEndpoint("common"),
		}
	)

	return microsoftOauthConfig, redirectURL, nil
}

func GenerateOTP() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func SendEmail(emailTo string, otp string) error {
	if err := godotenv.Load(); err != nil {
		return err
	}

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")
	
	msg := fmt.Sprintf("Subject: Your OTP Code\n\nYour OTP code is: %s", otp)
	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, []string{emailTo}, []byte(msg))
	return err
}
