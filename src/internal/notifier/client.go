package notifier

import "context"

type Client interface {
	NotifyToUser(ctx context.Context, message string, userID int64) error
}
