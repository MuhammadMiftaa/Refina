package helper

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/smtp"
	"time"

	"server/config/env"
	htmlTemplate "server/template"
)

type SMTPInterface interface {
	GetAuth() smtp.Auth
	GetAddress() string
	GetUser() string
}

type gmailSMTP struct {
	Host     string
	Port     string
	User     string
	Password string
	Secure   bool
	Auth     bool
}

func NewGmailSMTP(config env.GSMTP) SMTPInterface {
	return &gmailSMTP{
		Host:     config.GSHost,
		Port:     config.GSPort,
		User:     config.GSUser,
		Password: config.GSPassword,
	}
}

func (g *gmailSMTP) GetAuth() smtp.Auth {
	return smtp.PlainAuth("", g.User, g.Password, g.Host)
}

func (g *gmailSMTP) GetAddress() string {
	return g.Host + ":" + g.Port
}

func (g *gmailSMTP) GetUser() string {
	return g.User
}

type zohoSMTP struct {
	Host     string
	Port     string
	User     string
	Password string
	Secure   string
	Auth     bool
}

func NewZohoSMTP(config env.ZSMTP) SMTPInterface {
	return &zohoSMTP{
		Host:     config.ZSHost,
		Port:     config.ZSPort,
		User:     config.ZSUser,
		Password: config.ZSPassword,
		Secure:   config.ZSSecure,
		Auth:     config.ZSAuth,
	}
}

func (z *zohoSMTP) GetAuth() smtp.Auth {
	return smtp.PlainAuth("", z.User, z.Password, z.Host)
}

func (z *zohoSMTP) GetAddress() string {
	return z.Host + ":" + z.Port
}

func (z *zohoSMTP) GetUser() string {
	return z.User
}

type smtpClient struct {
	SMTPInterface SMTPInterface
}

type SMTPClientInterface interface {
	SendSingleEmail(to string, subject string, htmlFile string, data any) error
}

func NewSMTPClient(smtpInterface SMTPInterface) SMTPClientInterface {
	return &smtpClient{
		SMTPInterface: smtpInterface,
	}
}

func getTemplate(htmlFile string) (t *template.Template, err error) {
	t, err = template.New(htmlFile).Funcs(template.FuncMap{
		"formatDateMY": func(data time.Time) string {
			return data.Format("January 2006")
		},
		"formatDateMDY": func(data time.Time) string {
			return data.Format("January 02, 2006")
		},
		"formatDateMDYT": func(data time.Time) string {
			return data.Format("January 02, 2006 at 03:04 PM")
		},
		"convertBToMB": func(size int64) string {
			return fmt.Sprintf("%.2f", float64(size)/1024/1024)
		},
	}).Parse(htmlTemplate.Template[htmlFile])
	if err != nil {
		return nil, err
	}

	return t, nil
}

func parseHTML(file string, data any) (string, error) {
	bufferhtml := bytes.Buffer{}
	t, err := getTemplate(file)
	if err != nil {
		log.Printf("[ERROR] Failed to get template: %s", err)
		return "", err
	}
	// proses excecute data yang di masukkan dalam template html
	err = t.Execute(&bufferhtml, data)
	if err != nil {
		return "", err
	}

	return bufferhtml.String(), nil
}

func (c *smtpClient) SendSingleEmail(to string, subject string, htmlFile string, data any) error {
	auth := c.SMTPInterface.GetAuth()
	addr := c.SMTPInterface.GetAddress()
	user := c.SMTPInterface.GetUser()

	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"
	htmlFile, err := parseHTML(htmlFile, data)
	if err != nil {
		return err
	}

	msg := []byte(subject + mime + htmlFile)

	return smtp.SendMail(addr, auth, user, []string{to}, msg)
}
