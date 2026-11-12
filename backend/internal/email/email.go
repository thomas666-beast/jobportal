package email

import (
	"crypto/tls"
	"fmt"
	"log/slog"
	"net"
	"net/smtp"
	"os"
	"strconv"
	"strings"
	"time"
)

// Sender abstracts email delivery so we can swap implementations.
type Sender interface {
	Send(to, subject, htmlBody, textBody string) error
}

// SMTPConfig holds everything needed to send email via SMTP.
type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Pass     string
	From     string
	FromName string
	UseTLS   bool
	UseAuth  bool
}

// NewSMTPConfig reads config from environment variables.
// If SMTP_HOST is empty, returns nil so callers can fall back to NoopSender.
func NewSMTPConfig() *SMTPConfig {
	host := os.Getenv("SMTP_HOST")
	if host == "" {
		return nil
	}

	port := 1025
	if v := os.Getenv("SMTP_PORT"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			port = n
		}
	}

	from := os.Getenv("SMTP_FROM")
	if from == "" {
		from = "noreply@localhost"
	}

	fromName := os.Getenv("SMTP_FROM_NAME")
	if fromName == "" {
		fromName = "Job Portal"
	}

	return &SMTPConfig{
		Host:     host,
		Port:     port,
		User:     os.Getenv("SMTP_USER"),
		Pass:     os.Getenv("SMTP_PASS"),
		From:     from,
		FromName: fromName,
		UseTLS:   os.Getenv("SMTP_TLS") == "true",
		UseAuth:  os.Getenv("SMTP_AUTH") == "true",
	}
}

// SMTPSender sends email over plain SMTP.
type SMTPSender struct {
	cfg *SMTPConfig
}

func NewSMTPSender(cfg *SMTPConfig) *SMTPSender {
	return &SMTPSender{cfg: cfg}
}

func (s *SMTPSender) Send(to, subject, htmlBody, textBody string) error {
	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", s.cfg.FromName, s.cfg.From)
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"

	var body strings.Builder
	if htmlBody != "" {
		boundary := "----=_Part_" + strconv.FormatInt(time.Now().UnixNano(), 10)
		headers["Content-Type"] = "multipart/alternative; boundary=\"" + boundary + "\""

		body.WriteString("\r\n")
		body.WriteString("--" + boundary + "\r\n")
		body.WriteString("Content-Type: text/plain; charset=UTF-8\r\n\r\n")
		body.WriteString(textBody)
		body.WriteString("\r\n\r\n--" + boundary + "\r\n")
		body.WriteString("Content-Type: text/html; charset=UTF-8\r\n\r\n")
		body.WriteString(htmlBody)
		body.WriteString("\r\n\r\n--" + boundary + "--\r\n")
	} else {
		headers["Content-Type"] = "text/plain; charset=UTF-8"
		body.WriteString("\r\n")
		body.WriteString(textBody)
	}

	var msg strings.Builder
	for k, v := range headers {
		msg.WriteString(k + ": " + v + "\r\n")
	}
	msg.WriteString(body.String())

	addr := net.JoinHostPort(s.cfg.Host, strconv.Itoa(s.cfg.Port))

	var auth smtp.Auth
	if s.cfg.UseAuth && s.cfg.User != "" {
		auth = smtp.PlainAuth("", s.cfg.User, s.cfg.Pass, s.cfg.Host)
	}

	if s.cfg.UseTLS {
		return s.sendTLS(addr, auth, to, []byte(msg.String()))
	}
	return smtp.SendMail(addr, auth, s.cfg.From, []string{to}, []byte(msg.String()))
}

func (s *SMTPSender) sendTLS(addr string, auth smtp.Auth, to string, msg []byte) error {
	host, _, _ := net.SplitHostPort(addr)
	conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: host})
	if err != nil {
		return fmt.Errorf("tls dial: %w", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return fmt.Errorf("smtp client: %w", err)
	}
	defer client.Close()

	if auth != nil {
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("smtp auth: %w", err)
		}
	}
	if err := client.Mail(s.cfg.From); err != nil {
		return fmt.Errorf("smtp MAIL: %w", err)
	}
	if err := client.Rcpt(to); err != nil {
		return fmt.Errorf("smtp RCPT: %w", err)
	}
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("smtp DATA: %w", err)
	}
	if _, err := w.Write(msg); err != nil {
		return fmt.Errorf("smtp write: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("smtp close: %w", err)
	}
	return client.Quit()
}

// NoopSender logs instead of sending. Used when SMTP is not configured.
type NoopSender struct{}

func (NoopSender) Send(to, subject, htmlBody, textBody string) error {
	slog.Info("email suppressed (no-op)", "to", to, "subject", subject)
	return nil
}

// NewSender returns the appropriate sender based on env vars.
func NewSender() Sender {
	cfg := NewSMTPConfig()
	if cfg == nil {
		slog.Info("SMTP not configured — using no-op email sender")
		return NoopSender{}
	}
	slog.Info("SMTP configured", "host", cfg.Host, "port", cfg.Port, "tls", cfg.UseTLS, "auth", cfg.UseAuth)
	return NewSMTPSender(cfg)
}
