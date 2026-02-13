package notify

import (
	"context"
	"fmt"

	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"

	twilio "github.com/twilio/twilio-go"
)

type TwilioNotifier struct {
	client     *twilio.RestClient
	fromNumber string
}

func NewTwilioNotifier(accountSID, authToken, fromNumber string) *TwilioNotifier {
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})
	return &TwilioNotifier{
		client:     client,
		fromNumber: fromNumber,
	}
}

func (n *TwilioNotifier) Send(ctx context.Context, to string, _ string, body string) error {
	params := &twilioApi.CreateMessageParams{}
	params.SetTo(to)
	params.SetFrom(n.fromNumber)
	params.SetBody(body)

	_, err := n.client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("sending SMS via Twilio: %w", err)
	}
	return nil
}
