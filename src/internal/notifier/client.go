package notifier

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"gopkg.in/mail.v2"
)

type Client interface {
	NotifyToUser(ctx context.Context, message string, email string) error
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

func GetOptions() Options {
	host := os.Getenv("NOTIFIER_HOST")
	if host == "" {
		panic("notifier host is empty")
	}
	portStr := os.Getenv("NOTIFIER_PORT")
	if portStr == "" {
		panic("notifier port is empty")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		panic("notifier port is not a number")
	}
	sender := os.Getenv("NOTIFIER_SENDER")
	if sender == "" {
		panic("notifier sender is empty")
	}
	password := os.Getenv("NOTIFIER_PASSWORD")
	if password == "" {
		panic("notifier password is empty")
	}

	return Options{
		Host:     host,
		Port:     port,
		Username: sender,
		Password: password,
	}
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

func (c client) NotifyToUser(_ context.Context, message string, email string) error {
	m := mail.NewMessage()
	m.SetHeader("From", c.sender)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Transaction resume")
	m.SetBody("text/html", message)

	if err := c.dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("unexpected error sending mail to user due to: %w", err)
	}
	return nil
}
