package services

import "context"

type Notifier interface {
	Send(ctx context.Context, to string, subject string, body string) error
}
