package helper

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/smtp"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
	"unicode"

	"server/internal/types/dto"
	"server/internal/types/entity"
	htmlTemplate "server/template"

	"github.com/Rhymond/go-money"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/microsoft"
)

var mode = os.Getenv("MODE")

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
		return dto.UsersResponse{
			ID:    v.ID.String(),
			Name:  v.Name,
			Email: v.Email,
		}
	case entity.Transactions:
		return dto.TransactionsResponse{
			ID:              v.ID.String(),
			WalletID:        v.WalletID.String(),
			CategoryID:      v.CategoryID.String(),
			Amount:          v.Amount,
			TransactionDate: v.TransactionDate,
			Description:     v.Description,
		}
	case entity.Wallets:
		return dto.WalletsResponse{
			ID:           v.ID.String(),
			UserID:       v.UserID.String(),
			WalletTypeID: v.WalletTypeID.String(),
			Name:         v.Name,
			Number:       v.Number,
			Balance:      v.Balance,
		}
	case entity.Investments:
		return dto.InvestmentsResponse{
			ID:               v.ID.String(),
			UserID:           v.UserID.String(),
			InvestmentTypeID: v.InvestmentTypeID.String(),
			Name:             v.Name,
			Amount:           v.Amount,
			Quantity:         v.Quantity,
			InvestmentDate:   v.InvestmentDate,
			Description:      v.Description,
		}
	case entity.WalletTypes:
		return dto.WalletTypesResponse{
			ID:          v.ID.String(),
			Name:        v.Name,
			Type:        dto.WalletType(v.Type),
			Description: v.Description,
		}
	default:
		return nil
	}
}

var secretKey = "pojq09720ef1ko0f1h9iego2010j20240"

func GenerateToken(ID string, username string, email string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"id":       ID,
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

func VerifyToken(jwtToken string) (dto.UserData, error) {
	token, _ := jwt.Parse(jwtToken, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("parsing token error occured")
		}
		return []byte(secretKey), nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return dto.UserData{}, errors.New("token is invalid")
	}

	return dto.UserData{
		ID:       claims["id"].(string),
		Username: claims["username"].(string),
		Email:    claims["email"].(string),
	}, nil
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
	var (
		ClientID          = os.Getenv("GOOGLE_CLIENT_ID")
		ClientSecret      = os.Getenv("GOOGLE_CLIENT_SECRET")
		redirectURL       = os.Getenv("FRONTEND_URL")
		port              = os.Getenv("PORT")
		googleOauthConfig = &oauth2.Config{
			ClientID:     ClientID,
			ClientSecret: ClientSecret,
			RedirectURL:  "http://localhost:" + port + "/v1/auth/callback/google",
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		}
	)

	if mode == "production" {
		googleOauthConfig.RedirectURL = os.Getenv("PUBLIC_URL") + "/v1/auth/callback/google"
	}

	return googleOauthConfig, redirectURL, nil
}

func GetGithubOAuthConfig() (*oauth2.Config, string, error) {
	var (
		ClientID          = os.Getenv("GITHUB_CLIENT_ID")
		ClientSecret      = os.Getenv("GITHUB_CLIENT_SECRET")
		redirectURL       = os.Getenv("FRONTEND_URL")
		port              = os.Getenv("PORT")
		githubOauthConfig = &oauth2.Config{
			ClientID:     ClientID,
			ClientSecret: ClientSecret,
			RedirectURL:  "http://localhost:" + port + "/v1/auth/callback/github",
			Scopes: []string{
				"read:user",
				"user:email",
			},
			Endpoint: github.Endpoint,
		}
	)

	if mode == "production" {
		githubOauthConfig.RedirectURL = os.Getenv("PUBLIC_URL") + "/v1/auth/callback/github"
	}

	return githubOauthConfig, redirectURL, nil
}

func GetMicrosoftOAuthConfig() (*oauth2.Config, string, error) {
	var (
		ClientID             = os.Getenv("MICROSOFT_CLIENT_ID")
		ClientSecret         = os.Getenv("MICROSOFT_CLIENT_SECRET")
		redirectURL          = os.Getenv("FRONTEND_URL")
		port                 = os.Getenv("PORT")
		microsoftOauthConfig = &oauth2.Config{
			ClientID:     ClientID,
			ClientSecret: ClientSecret,
			RedirectURL:  "http://localhost:" + port + "/v1/auth/callback/microsoft",
			Scopes: []string{
				"User.Read",
			},
			Endpoint: microsoft.AzureADEndpoint("common"),
		}
	)

	if mode == "production" {
		microsoftOauthConfig.RedirectURL = os.Getenv("PUBLIC_URL") + "/v1/auth/callback/microsoft"
	}

	return microsoftOauthConfig, redirectURL, nil
}

func FormatAmountCurrency(amount int) string {
	format := money.NewFormatter(0, ".", ".", "RP", "RP 1")
	return format.Format(int64(amount))
}

func GetTemplate(htmlFile string) (t *template.Template, err error) {
	t, err = template.New(htmlFile).Funcs(template.FuncMap{
		"floatToString": func(data float64) string {
			return fmt.Sprintf("%0.f", data)
		},
		"formatCurency": func(data int) string {
			return FormatAmountCurrency(data)
		},
		"checkIsNotModZero": func(number int) bool {
			return number%2 != 0
		},
	}).Parse(htmlTemplate.OtpEmailTemplate)
	if err != nil {
		return nil, err
	}

	return t, nil
}

func ParseHTML(file string, otp string) (string, error) {
	bufferhtml := bytes.Buffer{}
	t, err := GetTemplate(file)
	if err != nil {
		log.Printf("[ERROR] Failed to get template: %s", err)
		return "", err
	}
	// proses excecute data yang di masukkan dalam template html
	err = t.Execute(&bufferhtml, struct {
		OTP string
	}{
		OTP: otp,
	})
	if err != nil {
		return "", err
	}

	return bufferhtml.String(), nil
}

func GenerateOTP() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func SendEmail(emailTo []string, otp string, htmlFilename string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")
	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)

	if htmlFilename == "" {
		return fmt.Errorf("html filename cannot be empty")
	}
	subject := "Subject: Welcome to Refina!\r\n"
	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"
	htmlFile, err := ParseHTML(htmlFilename, otp)
	if err != nil {
		return fmt.Errorf("failed to parse HTML template: %w", err)
	}

	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, smtpUser, emailTo, []byte(subject+mime+htmlFile))
	return err
}

func ParseUUID(id string) (uuid.UUID, error) {
	parsedID, err := uuid.Parse(id)
	if err != nil {
		return uuid.UUID{}, err
	}
	return parsedID, nil
}

func ExpandPathAndCreateDir(path string) (string, error) {
	// Ekspansi ~
	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		// Gabungkan home + sisanya
		path = filepath.Join(home, strings.TrimPrefix(path, "~"))
	}

	// Konversi ke path absolut
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	// Buat direktori jika belum ada
	err = os.MkdirAll(absPath, 0o755) // 0755 = rwxr-xr-x
	if err != nil {
		return "", err
	}

	return absPath, nil
}

func GenerateFileName(prefix, id string, postfix string) string {
	t := time.Now()
	timestamp := t.Format("20060102150405000000000")

	filename := fmt.Sprintf("%s_%s_%s_%s", prefix, id, timestamp, postfix)

	return filename
}
