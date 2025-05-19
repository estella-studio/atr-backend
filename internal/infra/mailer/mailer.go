package mailer

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/estella-studio/leon-backend/internal/infra/env"
	"github.com/google/uuid"
	"gopkg.in/gomail.v2"
)

type MailerItf interface {
	NewMail(to string, subject string, body string) error
}

type Mailer struct {
	Config *env.Env
}

type To struct {
	Email string `json:"email"`
}

type PasswordResetJSON struct {
	From struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"from"`
	To                []To      `json:"to"`
	TemplateUUID      uuid.UUID `json:"template_uuid"`
	TemplateVariables struct {
		UUID               uuid.UUID `json:"uuid"`
		CompanyInfoName    string    `json:"company_info_name"`
		CompanyInfoAddress string    `json:"company_info_address"`
		CompanyInfoCity    string    `json:"company_info_city"`
		CompanyInfoZipCode string    `json:"company_info_zip_code"`
		CompanyInfoCountry string    `json:"company_info_country"`
	} `json:"template_variables"`
}

func NewMailer(env *env.Env) *Mailer {
	return &Mailer{
		Config: env,
	}
}

func (m *Mailer) NewMail(to string, subject string, body string) error {
	mail := gomail.NewMessage()

	mail.SetHeader("From", m.Config.SMTPFrom)
	mail.SetHeader("To", to)
	mail.SetHeader("Subject", subject)
	mail.SetBody("text/html", body)

	d := gomail.NewDialer(
		m.Config.SMTPServer,
		m.Config.SMTPPort,
		m.Config.SMTPUsername,
		m.Config.SMTPPassword,
	)

	err := d.DialAndSend(mail)

	return err
}

func (m *Mailer) PasswordReset(to string, id uuid.UUID) error {
	url := "https://send.api.mailtrap.io/api/send"
	method := "POST"

	payload := &PasswordResetJSON{}

	payload.From.Email = m.Config.SMTPFrom
	payload.From.Name = "Estella Studio"
	payload.To = []To{{Email: to}}
	payload.TemplateUUID, _ = uuid.Parse("736645ed-1c8a-48d3-8e7a-388c32fd182e")
	payload.TemplateVariables.UUID = id
	payload.TemplateVariables.CompanyInfoName = "Estella Studio"
	payload.TemplateVariables.CompanyInfoAddress = "Malang"
	payload.TemplateVariables.CompanyInfoCity = "Malang"
	payload.TemplateVariables.CompanyInfoZipCode = "65144"
	payload.TemplateVariables.CompanyInfoCountry = "Indonesia"

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
	}

	payloadParsed := strings.NewReader(string(jsonBody))

	client := &http.Client{}

	req, err := http.NewRequest(method, url, payloadParsed)
	if err != nil {
		log.Println(err)
	}

	req.Header.Add(
		"Authorization",
		fmt.Sprintf("Bearer %s", "c8b4eae0c41260d8a14f5cee8bda744e"),
	)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}

	err = res.Body.Close()
	if err != nil {
		log.Println(err)
	}

	log.Println(string(body))

	return nil
}
