package mail

import (
	"fmt"
	"log"
	"net/mail"
	"net/smtp"
	"net/url"

	"crdx.org/dex/cmd/dexd/env"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
	FromName string
	To       string
}

func parseDSN(dsn string) (*Config, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return nil, fmt.Errorf("invalid DSN: %w", err)
	}

	if u.Scheme != "smtp" {
		return nil, fmt.Errorf("invalid scheme: expected smtp, got %s", u.Scheme)
	}

	config := &Config{
		Host: u.Hostname(),
		Port: u.Port(),
	}

	if config.Port == "" {
		config.Port = "587"
	}

	if u.User != nil {
		config.Username = u.User.Username()
		config.Password, _ = u.User.Password()
	}

	query := u.Query()
	config.From = query.Get("from")
	config.FromName = query.Get("from_name")
	config.To = query.Get("to")

	if config.From == "" {
		return nil, fmt.Errorf("missing 'from' parameter in DSN")
	}

	if config.To == "" {
		return nil, fmt.Errorf("missing 'to' parameter in DSN")
	}

	return config, nil
}

func send(config *Config, subject, body string) error {
	from := mail.Address{Name: config.FromName, Address: config.From}

	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s",
		from.String(),
		config.To,
		subject,
		body,
	)

	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)

	var auth smtp.Auth
	if config.Username != "" {
		auth = smtp.PlainAuth("", config.Username, config.Password, config.Host)
	}

	return smtp.SendMail(addr, auth, config.From, []string{config.To}, []byte(msg))
}

func Send(subject, body string) {
	dsn := env.NotifyDSN()
	if dsn == "" {
		return
	}

	config, err := parseDSN(dsn)
	if err != nil {
		log.Printf("mail: failed to parse DSN: %v", err)
		return
	}

	err = send(config, subject, body)
	if err != nil {
		log.Printf("mail: failed to send email: %v", err)
	}
}
