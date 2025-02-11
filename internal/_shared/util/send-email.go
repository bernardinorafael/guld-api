package util

import (
	"github.com/resend/resend-go/v2"
)

type EmailClient struct {
	client *resend.Client
}

type EmailBody struct {
	Subject string
	Body    string
	To      string
}

func NewEmailClient(resendApiKey string) *EmailClient {
	return &EmailClient{
		client: resend.NewClient(resendApiKey),
	}
}

func (e *EmailClient) SendEmail(input EmailBody) error {
	params := &resend.SendEmailRequest{
		From:    "Notification <onboarding@resend.dev>",
		To:      []string{input.To},
		Html:    input.Body,
		Subject: input.Subject,
	}

	_, err := e.client.Emails.Send(params)
	if err != nil {
		return err
	}

	return nil
}
