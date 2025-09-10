package notification

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"net/mail"
	"net/smtp"
	"strings"
)

// Mailer is a configurable mail sender that uses embedded templates for HTML emails.
type Mailer struct {
	host       string
	port       int
	username   string
	password   string
	fromName   string
	fromAddr   string
	templateFS embed.FS
	auth       smtp.Auth

	// logger for mailer internal logging (optional, default provided)
	logger *slog.Logger
}

// MailerOptions holds configuration for NewMailer.
type MailerOptions struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string

	FromName    string
	FromAddress string

	TemplateFS embed.FS

	// optional logger; if nil a default slog.Logger will be created
	Logger *slog.Logger
}

// sanitizeHeader trims whitespace/quotes and strips CR/LF to prevent header injection.
func sanitizeHeader(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Trim(s, "\"")
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\n", "")
	return s
}

// NewMailer creates a configured Mailer from MailerOptions.
// Required fields:
//   - SMTPHost
//   - FromAddress
//   - TemplateFS
//
// SMTPPort defaults to 587 when zero. If Logger is nil a default slog.Logger is created.
func NewMailer(opts MailerOptions) (*Mailer, error) {
	fromAddr := sanitizeHeader(opts.FromAddress)
	fromName := sanitizeHeader(opts.FromName)

	m := &Mailer{
		host:       opts.SMTPHost,
		port:       587,
		username:   opts.SMTPUsername,
		password:   opts.SMTPPassword,
		fromName:   fromName,
		fromAddr:   fromAddr,
		templateFS: opts.TemplateFS,
		logger:     opts.Logger,
	}

	if opts.SMTPPort != 0 {
		m.port = opts.SMTPPort
	}

	// validate required fields
	if m.host == "" {
		return nil, errors.New("smtp host is required")
	}
	if m.fromAddr == "" {
		return nil, errors.New("from address is required")
	}

	// validate email format
	if _, err := mail.ParseAddress(m.fromAddr); err != nil {
		return nil, errors.New("invalid from address")
	}

	// ensure fromName doesn't contain angle brackets (avoid header injection)
	m.fromName = strings.ReplaceAll(m.fromName, "<", "")
	m.fromName = strings.ReplaceAll(m.fromName, ">", "")

	// check templateFS was provided
	if m.templateFS == (embed.FS{}) {
		return nil, errors.New("template FS is required")
	}

	// ensure logger
	if m.logger == nil {
		m.logger = slog.New(slog.NewJSONHandler(nil, &slog.HandlerOptions{}))
	}

	if m.username != "" && m.password != "" {
		m.auth = smtp.PlainAuth("", m.username, m.password, m.host)
	}

	m.logger.Debug("mailer configured", "host", m.host, "port", m.port, "from", m.fromAddr)
	return m, nil
}

// SendEmail sends an HTML email rendered from an embedded template.
// templateName should match a file under the embedded "emails/" directory (e.g. "welcome.html").
// 'to' accepts one or more recipient addresses.
func (m *Mailer) SendEmail(ctx context.Context, to []string, subject, templateName string, data any) error {
	if len(to) == 0 {
		return errors.New("at least one recipient required")
	}

	body, err := m.formatEmailHTML(templateName, data)
	if err != nil {
		return err
	}

	headers := map[string]string{
		"From":         fmt.Sprintf("%s <%s>", m.fromName, m.fromAddr),
		"To":           strings.Join(to, ", "),
		"Subject":      subject,
		"MIME-Version": "1.0",
		"Content-Type": `text/html; charset="utf-8"`,
	}

	var msg bytes.Buffer
	for k, v := range headers {
		msg.WriteString(fmt.Sprintf("%s: %s\r\n", k, v))
	}
	msg.WriteString("\r\n")
	msg.WriteString(body)

	addr := fmt.Sprintf("%s:%d", m.host, m.port)
	var auth smtp.Auth
	if m.auth != nil {
		auth = m.auth
	} else {
		auth = nil
	}

	// net/smtp.SendMail will use STARTTLS if the server supports it.
	if err := smtp.SendMail(addr, auth, m.fromAddr, to, msg.Bytes()); err != nil {
		m.logger.Error("failed to send email", "err", err, "to", to, "subject", subject)
		return err
	}
	m.logger.Debug("email sent", "to", to, "subject", subject)
	return nil
}

// formatEmailHTML loads the named template from the embedded FS and executes it with data.
// Template files should be located under "emails/" in the embedded FS (see templates/embed.go).
func (m *Mailer) formatEmailHTML(templateName string, data any) (string, error) {
	if templateName == "" {
		return "", errors.New("template name required")
	}

	// Load the specific template file from the embedded FS (e.g. "emails/welcome.html").
	tplPath := "emails/" + templateName
	t, err := template.ParseFS(m.templateFS, tplPath)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
