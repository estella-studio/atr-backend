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
	AccountRegistration(to string, code uint) error
	PasswordReset(to string, id uuid.UUID, code uint) error
}

type Mailer struct {
	Config *env.Env
}

type To struct {
	Email string `json:"email"`
}

type PasswordReset struct {
	From struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"from"`
	To                []To      `json:"to"`
	TemplateUUID      uuid.UUID `json:"template_uuid"`
	TemplateVariables struct {
		UUID               uuid.UUID `json:"uuid"`
		Code               uint      `json:"code"`
		CompanyInfoName    string    `json:"company_info_name"`
		CompanyInfoAddress string    `json:"company_info_address"`
		CompanyInfoCity    string    `json:"company_info_city"`
		CompanyInfoZipCode string    `json:"company_info_zip_code"`
		CompanyInfoCountry string    `json:"company_info_country"`
	} `json:"template_variables"`
}

type AccountRegistration struct {
	From struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"from"`
	To                []To      `json:"to"`
	TemplateUUID      uuid.UUID `json:"template_uuid"`
	TemplateVariables struct {
		Code               uint   `json:"code"`
		CompanyInfoName    string `json:"company_info_name"`
		CompanyInfoAddress string `json:"company_info_address"`
		CompanyInfoCity    string `json:"company_info_city"`
		CompanyInfoZipCode string `json:"company_info_zip_code"`
		CompanyInfoCountry string `json:"company_info_country"`
	} `json:"template_variables"`
}

func NewMailer(env *env.Env) MailerItf {
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

func (m *Mailer) AccountRegistration(to string, code uint) error {
	url := m.Config.MailtrapURL
	method := "POST"

	payload := &AccountRegistration{}

	payload.From.Email = m.Config.SMTPFrom
	payload.From.Name = m.Config.EmailFrom
	payload.To = []To{{Email: to}}
	payload.TemplateUUID, _ = uuid.Parse(m.Config.MailtrapTemplateAccountRegistration)
	payload.TemplateVariables.Code = code
	payload.TemplateVariables.CompanyInfoName = m.Config.MailtrapCompanyInfoName
	payload.TemplateVariables.CompanyInfoAddress = m.Config.MailtrapCompanyInfoAddress
	payload.TemplateVariables.CompanyInfoCity = m.Config.MailtrapCompanyInfoCity
	payload.TemplateVariables.CompanyInfoZipCode = m.Config.MailtrapCompanyInfoZipCode
	payload.TemplateVariables.CompanyInfoCountry = m.Config.MailtrapCompanyInfoCountry

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
		fmt.Sprintf("Bearer %s", m.Config.MailtrapTokenAccountRegistration),
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

	return err
}

func (m *Mailer) PasswordReset(to string, id uuid.UUID, code uint) error {
	url := m.Config.MailtrapURL
	method := "POST"

	payload := &PasswordReset{}

	payload.From.Email = m.Config.SMTPFrom
	payload.From.Name = m.Config.EmailFrom
	payload.To = []To{{Email: to}}
	payload.TemplateUUID, _ = uuid.Parse(m.Config.MailtrapTemplatePasswordReset)
	payload.TemplateVariables.UUID = id
	payload.TemplateVariables.Code = code
	payload.TemplateVariables.CompanyInfoName = m.Config.MailtrapCompanyInfoName
	payload.TemplateVariables.CompanyInfoAddress = m.Config.MailtrapCompanyInfoAddress
	payload.TemplateVariables.CompanyInfoCity = m.Config.MailtrapCompanyInfoCity
	payload.TemplateVariables.CompanyInfoZipCode = m.Config.MailtrapCompanyInfoZipCode
	payload.TemplateVariables.CompanyInfoCountry = m.Config.MailtrapCompanyInfoCountry

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
		fmt.Sprintf("Bearer %s", m.Config.MailtrapTokenPasswordReset),
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

	return err
}
