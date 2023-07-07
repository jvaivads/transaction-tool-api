package summarizer

import (
	"context"
	"fmt"
	"transaction-tool-api/src/internal/notifier"
)

type Service interface {
	notifyResume(ctx context.Context, txns transactions) (err error)
}

type service struct {
	repository repository
	notifier   notifier.Client
}

func (s service) notifyResume(ctx context.Context, txns transactions) (err error) {
	if len(txns.items) == 0 {
		return nil
	}

	var (
		summ   summarizer
		repoTx tx
	)

	if repoTx, err = s.repository.initTransactionalOperations(ctx); err != nil {
		err = fmt.Errorf("error creating repository transaction due to: %w", err)
		return
	}
	defer func() {
		err = s.repository.finishTransactionalOperations(ctx, repoTx, err)
	}()

	if err = s.repository.saveBankTransactions(ctx, repoTx, txns); err != nil {
		err = fmt.Errorf("error saving transactions due to: %w", err)
		return
	}

	if err = s.notifier.NotifyToUser(ctx, summ.resume(txns).String(), txns.userID); err != nil {
		err = fmt.Errorf("error notifying transactions to user id %d due to: %w", txns.userID, err)
		return
	}

	return nil
}
