package mailer

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"time"

	"html/template"

	"github.com/bernardinorafael/pkg/logger"
	"github.com/resend/resend-go/v2"
)

//go:embed "templates"
var templateFS embed.FS

type SendParams struct {
	From    string
	To      string
	Subject string
	File    string
	Data    any
}

type Mailer struct {
	ctx    context.Context
	log    logger.Logger
	client *resend.Client
}

func New(ctx context.Context, log logger.Logger, apiKey string) Mailer {
	return Mailer{
		ctx:    ctx,
		log:    log,
		client: resend.NewClient(apiKey),
	}
}

func (m *Mailer) Send(p SendParams) error {
	tmpl, err := template.New("email").ParseFS(templateFS, "templates/"+p.File)
	if err != nil {
		m.log.Errorw(m.ctx, "error on parse template", logger.Err(err))
		return fmt.Errorf("error on parse template: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, p.Data); err != nil {
		m.log.Errorw(m.ctx, "error on execute template", logger.Err(err))
		return fmt.Errorf("error on execute template: %w", err)
	}

	params := &resend.SendEmailRequest{
		From:    p.From,
		To:      []string{p.To},
		Html:    body.String(),
		Subject: p.Subject,
	}

	return m.send(params, 3)
}

func (m *Mailer) send(params *resend.SendEmailRequest, maxRetries int) error {
	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		_, err := m.client.Emails.Send(params)
		if err == nil {
			m.log.Infow(m.ctx, "email sent", logger.Int("attempt", attempt))
			return nil
		}

		m.log.Errorw(
			m.ctx,
			"error on send email",
			logger.Err(err),
			logger.Int("attempt", attempt),
		)
		lastErr = err
		time.Sleep(time.Second * 2)
	}

	return fmt.Errorf("error on send email after %d attempts: %w", maxRetries, lastErr)
}
