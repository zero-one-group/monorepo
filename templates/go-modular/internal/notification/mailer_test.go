package notification

import (
	"context"
	"net/smtp"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMailer(t *testing.T) {
	// prepare a fake FS with an emails/welcome.html template
	mfs := fstest.MapFS{
		"emails/welcome.html": &fstest.MapFile{
			Data: []byte("<h1>Hello {{.Name}}</h1>"),
		},
	}

	t.Run("NewMailer_validation_errors", func(t *testing.T) {
		_, err := NewMailer(MailerOptions{
			SMTPHost:    "",
			FromAddress: "sender@example.com",
			TemplateFS:  mfs,
		})
		require.Error(t, err)

		_, err = NewMailer(MailerOptions{
			SMTPHost:    "smtp.example",
			FromAddress: "",
			TemplateFS:  mfs,
		})
		require.Error(t, err)

		_, err = NewMailer(MailerOptions{
			SMTPHost:    "smtp.example",
			FromAddress: "sender@example.com",
			TemplateFS:  nil,
		})
		require.Error(t, err)
	})

	t.Run("FormatEmailHTML_not_found", func(t *testing.T) {
		mailer, err := NewMailer(MailerOptions{
			SMTPHost:    "smtp.example",
			FromAddress: "sender@example.com",
			TemplateFS:  mfs,
		})
		require.NoError(t, err)

		_, err = mailer.formatEmailHTML("missing.html", nil)
		require.Error(t, err)
	})

	t.Run("SendEmail_success_calls_sendMail", func(t *testing.T) {
		mailer, err := NewMailer(MailerOptions{
			SMTPHost:    "smtp.test",
			FromAddress: "sender@example.com",
			TemplateFS:  mfs,
		})
		require.NoError(t, err)

		// override sendMail to capture args
		orig := sendMail
		defer func() { sendMail = orig }()

		var gotAddr string
		var gotFrom string
		var gotTo []string
		var gotMsg []byte

		sendMail = func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
			gotAddr = addr
			gotFrom = from
			gotTo = to
			gotMsg = msg
			return nil
		}

		err = mailer.SendEmail(context.Background(), []string{"to@example.com"}, "Welcome", "welcome.html", map[string]string{"Name": "Alice"})
		require.NoError(t, err)

		assert.Contains(t, gotAddr, "smtp.test")
		assert.Equal(t, "sender@example.com", gotFrom)
		require.Len(t, gotTo, 1)
		assert.Equal(t, "to@example.com", gotTo[0])
		assert.Contains(t, string(gotMsg), "Hello Alice")
	})
}
