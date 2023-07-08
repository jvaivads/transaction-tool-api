package notifier

import (
	"context"
	"fmt"

	"gopkg.in/mail.v2"
)

type Client interface {
	NotifyToUser(ctx context.Context, message string, userID int64) error
}

type dialer interface {
	DialAndSend(m ...*mail.Message) error
}

type client struct {
	sender string
	dialer dialer
}

type Options struct {
	Host     string
	Port     int
	Username string
	Password string
}

func NewClient(options Options) Client {
	return client{
		sender: options.Username,
		dialer: mail.NewDialer(
			options.Host,
			options.Port,
			options.Username,
			options.Password,
		),
	}
}

func (c client) NotifyToUser(_ context.Context, message string, userID int64) error {
	m := mail.NewMessage()
	m.SetHeader("From", c.sender)
	m.SetHeader("To", "to@example.com")
	m.SetHeader("Subject", "transaction resume")
	m.SetBody("text/plain", message)

	if err := c.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("unexpected error sending mail to usier id %d due to: %w", userID, err)
	}
	return nil
}
