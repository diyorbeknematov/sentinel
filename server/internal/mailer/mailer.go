package mailer

import (
	"fmt"
	"net/smtp"
)

type Mailer struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func New(host, port, username, password, from string) *Mailer {
	return &Mailer{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

func (m *Mailer) Send(to, subject, body string) error {
	auth := smtp.PlainAuth("", m.username, m.password, m.host)

	msg := fmt.Sprintf(
		"From: %s\r\nSubject: %s\r\n\r\n%s",
		m.from, subject, body,
	)

	addr := m.host + ":" + m.port

	return smtp.SendMail(addr, auth, m.from, []string{to}, []byte(msg))
}
