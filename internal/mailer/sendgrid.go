package mailer

import (
	"fmt"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/youssefM1999/report/pkg/retry"
)

type SendGridMailer struct {
	from   string
	apiKey string
	client *sendgrid.Client
}

func NewSendGridMailer(from, apiKey string) *SendGridMailer {
	client := sendgrid.NewSendClient(apiKey)
	return &SendGridMailer{
		from:   from,
		apiKey: apiKey,
		client: client,
	}
}

func (m *SendGridMailer) Send(email, username, subject, markdownContent string, period time.Duration, isSandbox bool) (int, error) {
	from := mail.NewEmail(FromName, m.from)
	to := mail.NewEmail(username, email)

	data := NewEmailData(subject, markdownContent, period)

	body, err := renderEmailTemplate(data)
	if err != nil {
		return -1, fmt.Errorf("failed to render email template: %w", err)
	}

	message := mail.NewSingleEmail(from, data.Subject, to, "", body)

	message.SetMailSettings(&mail.MailSettings{
		SandboxMode: &mail.Setting{
			Enable: &isSandbox,
		},
	})

	var statusCode int
	err = retry.Retry(func() error {
		response, err := m.client.Send(message)
		statusCode = response.StatusCode
		return err
	}, MaxRetries)
	if err != nil {
		return -1, fmt.Errorf("failed to send email: %w after %d retries", err, MaxRetries)
	}

	return statusCode, nil
}
