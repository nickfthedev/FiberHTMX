package lib

import (
	"crypto/tls"
	"net/smtp"
)

func SendEmail(recipient string, subject string, body string, contentType string) error {
	smtpHost := Config.SMTPHost
	smtpPort := Config.SMTPPort
	smtpUser := Config.SMTPUser
	smtpPass := Config.SMTPPass

	auth := smtp.PlainAuth("", smtpUser, smtpPass, smtpHost)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}

	// Connect to the SMTP Server
	c, err := smtp.Dial(smtpHost + ":" + smtpPort)

	if err != nil {
		return err
	}

	// start TLS for an encrypted connection
	if err = c.StartTLS(tlsconfig); err != nil {
		return err
	}

	// Auth
	if err = c.Auth(auth); err != nil {
		return err
	}

	// To && From
	if err = c.Mail(smtpUser); err != nil {
		return err
	}
	if err = c.Rcpt(recipient); err != nil {
		return err
	}

	// Data
	w, err := c.Data()
	if err != nil {
		return err
	}

	msg := "From: " + smtpUser + "\n" +
		"To: " + recipient + "\n" +
		"Subject: " + subject + "\n" +
		"MIME-version: 1.0;\nContent-Type: " + contentType + "; charset=\"UTF-8\";\n\n" +
		body

	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	c.Quit()

	return nil
}
