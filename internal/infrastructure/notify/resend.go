package notify

import (
	"context"
	"fmt"

	"github.com/resend/resend-go/v2"
)

type ResendNotifier struct {
	client *resend.Client
	from   string
}

func NewResendNotifier(apiKey, from string) *ResendNotifier {
	return &ResendNotifier{
		client: resend.NewClient(apiKey),
		from:   from,
	}
}

func (n *ResendNotifier) Send(ctx context.Context, to string, subject string, body string) error {
	params := &resend.SendEmailRequest{
		From:    n.from,
		To:      []string{to},
		Subject: subject,
		Text:    body,
	}
	_, err := n.client.Emails.SendWithContext(ctx, params)
	if err != nil {
		return fmt.Errorf("sending email via Resend: %w", err)
	}
	return nil
}
