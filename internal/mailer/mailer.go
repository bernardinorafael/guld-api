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

type Config struct {
	APIKey           string
	MaxRetries       int
	RetryDelay       time.Duration
	OperationTimeout time.Duration
}

type Mailer struct {
	ctx    context.Context
	log    logger.Logger
	client *resend.Client
	config Config
}

func New(ctx context.Context, log logger.Logger, config Config) Mailer {
	// NOTE: Verify default values
	if config.MaxRetries == 0 {
		config.MaxRetries = 3
	}
	// NOTE: Verify default retry delay
	if config.RetryDelay == 0 {
		config.RetryDelay = time.Second * 2
	}
	// NOTE: Verify default operation timeout
	if config.OperationTimeout == 0 {
		config.OperationTimeout = time.Second * 10
	}

	return Mailer{
		ctx:    ctx,
		log:    log,
		client: resend.NewClient(config.APIKey),
		config: config,
	}
}

func (m *Mailer) Send(p SendParams) error {
	if err := m.validateParams(p); err != nil {
		return fmt.Errorf("invalid parameters: %w", err)
	}

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

	return m.send(params, m.config.MaxRetries)
}

func (m *Mailer) validateParams(p SendParams) error {
	if p.From == "" {
		return fmt.Errorf("from email is required")
	}
	if p.To == "" {
		return fmt.Errorf("to email is required")
	}
	if p.Subject == "" {
		return fmt.Errorf("subject is required")
	}
	if p.File == "" {
		return fmt.Errorf("template file is required")
	}
	return nil
}

func (m *Mailer) send(params *resend.SendEmailRequest, maxRetries int) error {
	var lastErr error
	for attempt := 1; attempt <= maxRetries; attempt++ {
		ctx, cancel := context.WithTimeout(m.ctx, m.config.OperationTimeout)
		defer cancel()

		_, err := m.client.Emails.Send(params)
		if err == nil {
			m.log.Infow(ctx, "email sent",
				logger.Int("attempt", attempt),
				logger.String("to", params.To[0]),
			)
			return nil
		}

		m.log.Errorw(
			ctx,
			"error on send email",
			logger.Err(err),
			logger.Int("attempt", attempt),
			logger.String("to", params.To[0]),
		)
		lastErr = err
		time.Sleep(m.config.RetryDelay)
	}

	return fmt.Errorf("error on send email after %d attempts: %w", maxRetries, lastErr)
}
