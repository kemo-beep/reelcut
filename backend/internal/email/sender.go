package email

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/smtp"
	"strings"
)

// Sender sends transactional email (e.g. password reset, verification).
type Sender interface {
	Send(to, subject, body string) error
}

// NoOpSender logs the email and does not send (e.g. when no SMTP config).
type NoOpSender struct{}

func (NoOpSender) Send(to, subject, body string) error {
	log.Printf("[email] no-op send to=%q subject=%q (set SMTP_* or SENDGRID_API_KEY to send)", to, subject)
	return nil
}

// SMTPConfig for sending via any SMTP server (including SendGrid SMTP relay).
type SMTPConfig struct {
	From     string // e.g. noreply@example.com
	Host     string // e.g. smtp.sendgrid.net
	Port     int    // e.g. 587
	Username string // e.g. apikey for SendGrid
	Password string
	UseTLS   bool
}

// SMTPSender sends email via SMTP.
type SMTPSender struct {
	cfg SMTPConfig
}

func NewSMTPSender(cfg SMTPConfig) *SMTPSender {
	return &SMTPSender{cfg: cfg}
}

func (s *SMTPSender) Send(to, subject, body string) error {
	if s.cfg.Host == "" || s.cfg.Port == 0 {
		log.Printf("[email] SMTP not configured, skipping send to=%q subject=%q", to, subject)
		return nil
	}
	addr := fmt.Sprintf("%s:%d", s.cfg.Host, s.cfg.Port)
	msg := []byte(
		"From: " + s.cfg.From + "\r\n" +
			"To: " + to + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"Content-Type: text/plain; charset=UTF-8\r\n" +
			"\r\n" + body + "\r\n",
	)
	auth := smtp.PlainAuth("", s.cfg.Username, s.cfg.Password, s.cfg.Host)
	if s.cfg.UseTLS {
		return s.sendTLS(addr, to, msg, auth)
	}
	return smtp.SendMail(addr, auth, s.cfg.From, []string{to}, msg)
}

func (s *SMTPSender) sendTLS(addr, to string, msg []byte, auth smtp.Auth) error {
	host := strings.Split(addr, ":")[0]
	conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: host})
	if err != nil {
		return err
	}
	defer conn.Close()
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer client.Close()
	if err = client.Auth(auth); err != nil {
		return err
	}
	if err = client.Mail(s.cfg.From); err != nil {
		return err
	}
	if err = client.Rcpt(to); err != nil {
		return err
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return client.Quit()
}
