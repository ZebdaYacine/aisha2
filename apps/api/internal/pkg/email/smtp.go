package email

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/mail"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

var ErrNotConfigured = errors.New("SMTP email is not configured")

type Message struct {
	To      string
	Subject string
	Body    string
}

type SMTP struct {
	host     string
	port     int
	user     string
	password string
	from     string
}

func NewSMTP(host string, port int, user, password, from string) *SMTP {
	if strings.TrimSpace(host) == "" || strings.TrimSpace(user) == "" || strings.TrimSpace(password) == "" || strings.TrimSpace(from) == "" {
		return nil
	}
	return &SMTP{host: host, port: port, user: user, password: password, from: from}
}

func (s *SMTP) Send(ctx context.Context, to, subject, body string) error {
	return s.SendMessage(ctx, Message{To: to, Subject: subject, Body: body})
}

func (s *SMTP) SendMessage(ctx context.Context, message Message) error {
	if s == nil {
		return ErrNotConfigured
	}
	from, err := mail.ParseAddress(s.from)
	if err != nil || from.Address == "" {
		return fmt.Errorf("invalid SMTP_FROM: %w", err)
	}
	to, err := mail.ParseAddress(message.To)
	if err != nil || to.Address == "" {
		return fmt.Errorf("invalid recipient address: %w", err)
	}
	if strings.ContainsAny(message.Subject, "\r\n") {
		return errors.New("email subject contains invalid line breaks")
	}

	port := s.port
	if port < 1 || port > 65535 {
		return fmt.Errorf("invalid SMTP port: %d", port)
	}
	address := net.JoinHostPort(s.host, strconv.Itoa(port))
	dialer := &net.Dialer{}
	connection, err := dialer.DialContext(ctx, "tcp", address)
	if err != nil {
		return fmt.Errorf("connect to SMTP server: %w", err)
	}
	defer connection.Close()
	_ = connection.SetDeadline(time.Now().Add(30 * time.Second))

	var client *smtp.Client
	if port == 465 {
		tlsConnection := tls.Client(connection, &tls.Config{ServerName: s.host, MinVersion: tls.VersionTLS12})
		if err := tlsConnection.HandshakeContext(ctx); err != nil {
			return fmt.Errorf("SMTP TLS handshake: %w", err)
		}
		client, err = smtp.NewClient(tlsConnection, s.host)
	} else {
		client, err = smtp.NewClient(connection, s.host)
		if err == nil {
			if ok, _ := client.Extension("STARTTLS"); ok {
				err = client.StartTLS(&tls.Config{ServerName: s.host, MinVersion: tls.VersionTLS12})
			}
		}
	}
	if err != nil {
		return fmt.Errorf("create SMTP client: %w", err)
	}
	defer client.Close()

	if err := client.Auth(smtp.PlainAuth("", s.user, s.password, s.host)); err != nil {
		return fmt.Errorf("authenticate with SMTP server: %w", err)
	}
	if err := client.Mail(from.Address); err != nil {
		return fmt.Errorf("set SMTP sender: %w", err)
	}
	if err := client.Rcpt(to.Address); err != nil {
		return fmt.Errorf("set SMTP recipient: %w", err)
	}
	writer, err := client.Data()
	if err != nil {
		return fmt.Errorf("open SMTP message: %w", err)
	}
	content := "From: " + from.String() + "\r\n" +
		"To: " + to.String() + "\r\n" +
		"Subject: " + message.Subject + "\r\n" +
		"MIME-Version: 1.0\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"Content-Transfer-Encoding: 8bit\r\n\r\n" +
		message.Body
	if _, err := io.WriteString(writer, content); err != nil {
		_ = writer.Close()
		return fmt.Errorf("write SMTP message: %w", err)
	}
	if err := writer.Close(); err != nil {
		return fmt.Errorf("finish SMTP message: %w", err)
	}
	if err := client.Quit(); err != nil {
		return fmt.Errorf("close SMTP session: %w", err)
	}
	return nil
}
