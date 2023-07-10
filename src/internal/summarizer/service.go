package summarizer

import (
	"context"
	"fmt"
	"transaction-tool-api/src/internal/notifier"
)

func NewService(repository Repository, notifierClient notifier.Client) Service {
	return service{
		repository: repository,
		notifier:   notifierClient,
	}
}

type Service interface {
	notifyResume(ctx context.Context, txns transactions) (err error)
}

type service struct {
	repository Repository
	notifier   notifier.Client
}

func (s service) notifyResume(ctx context.Context, txns transactions) (err error) {
	if len(txns.items) == 0 {
		return nil
	}

	var (
		summ    summarizer
		user    User
		repoTx  tx
		message string
	)

	if repoTx, err = s.repository.initTransactionalOperations(ctx); err != nil {
		err = fmt.Errorf("error creating repository transaction due to: %w", err)
		return
	}
	defer func() {
		err = s.repository.finishTransactionalOperations(ctx, repoTx, err)
	}()

	if user, err = s.repository.getUserByID(ctx, repoTx, txns.userID); err != nil {
		err = fmt.Errorf("error getting user due to: %w", err)
		return
	}

	if err = s.repository.saveBankTransactions(ctx, repoTx, txns); err != nil {
		err = fmt.Errorf("error saving transactions due to: %w", err)
		return
	}

	if message, err = summ.resume(user, txns).ToHTML(resumeHTMLTemplate); err != nil {
		err = fmt.Errorf("error generating message for user id %d due to: %w", txns.userID, err)
		return
	}

	if err = s.notifier.NotifyToUser(ctx, message, user.Email); err != nil {
		err = fmt.Errorf("error notifying transactions to user id %d due to: %w", txns.userID, err)
		return
	}

	return nil
}
