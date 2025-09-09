package email

import (
	"bytes"
	"fmt"
	"html/template"
	"net/smtp"

	"github.com/AlexMickh/coledzh-shop-backend/internal/config"
)

type VerificationEmailVars struct {
	Login string
	Token string
}

type Email struct {
	cfg  config.MailConfig
	auth smtp.Auth
}

func New(cfg config.MailConfig) *Email {
	return &Email{
		cfg:  cfg,
		auth: smtp.PlainAuth("", cfg.FromAddr, cfg.Password, cfg.Host),
	}
}

func (e *Email) SendVerification(to string, token, login string) error {
	const op = "lib.email.Send"

	tmpl, err := template.ParseFiles("./internal/lib/email/templates/verify-email.html")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rendered := new(bytes.Buffer)
	vars := VerificationEmailVars{
		Login: login,
		Token: token,
	}
	if err = tmpl.Execute(rendered, vars); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"

	err = smtp.SendMail(
		fmt.Sprintf("%s:%d", e.cfg.Host, e.cfg.Port),
		e.auth,
		e.cfg.FromAddr,
		[]string{to},
		fmt.Appendf(nil, "Subject: Email\n%s\n\n%s", headers, rendered.String()),
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
