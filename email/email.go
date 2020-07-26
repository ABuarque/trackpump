package email

import (
	"fmt"
	"net/smtp"
)

// Client defines the state of this object
type Client struct {
	Email    string
	Password string
}

// New returns a new email instance
func New(email, password string) *Client {
	return &Client{
		Email:    email,
		Password: password,
	}
}

// Send sends an email
func (e *Client) Send(from string, to []string, subject, body string) error {
	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", e.Email, e.Password, "smtp.gmail.com"),
		e.Email, to, []byte(body))
	if err != nil {
		return fmt.Errorf("falha ao enviar email, erro %q", err)
	}
	return nil
}
