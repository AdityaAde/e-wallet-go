package service

import (
	"net/smtp"

	"adityaad.id/belajar-auth/domain"
	"adityaad.id/belajar-auth/internal/config"
)

type emailService struct {
	cnf *config.Config
}

func NewEmail(cnf *config.Config) domain.EmailService {
	return &emailService{
		cnf: cnf,
	}
}

// Send implements domain.EmailService.
func (e emailService) Send(to string, subject string, body string) error {
	auth := smtp.PlainAuth(
		"",
		e.cnf.Mail.User,
		e.cnf.Mail.Password,
		e.cnf.Mail.Host,
	)

	msg := []byte("" + "From: Aditya <" + e.cnf.Mail.User + ">\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" + body)

	return smtp.SendMail(
		e.cnf.Mail.Host+":"+e.cnf.Mail.Port,
		auth,
		e.cnf.Mail.User,
		[]string{to},
		msg,
	)
}
