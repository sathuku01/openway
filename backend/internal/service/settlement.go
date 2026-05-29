package service

import (
	"context"
	"errors"

	"openway/internal/queue"
	"openway/internal/repository"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type SettlementService struct {
	repo        *repository.SettlementRepository
	queueClient *asynq.Client
}

func NewSettlementService(repo *repository.SettlementRepository, client *asynq.Client) *SettlementService {
	return &SettlementService{repo: repo, queueClient: client}
}

func (s *SettlementService) ProcessWebhook(ctx context.Context, telcoRef, sender, receiver, currency string, amount float64, rawPayload []byte) error {
	settlement := &repository.Settlement{
		TelcoReference: telcoRef,
		SenderPhone:    sender,
		ReceiverPhone:  receiver,
		Amount:         amount,
		Currency:       currency,
		RawPayload:     rawPayload,
		Status:         repository.StatusPending,
	}

	err := s.repo.Create(ctx, settlement)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateReference) {
			log.Warn().Str("ref", telcoRef).Msg("duplicate webhook received, ignoring")
			return nil
		}
		return err
	}

	payload := queue.SettlementProcessPayload{
		SettlementID:   settlement.ID,
		TelcoReference: settlement.TelcoReference,
		Amount:         settlement.Amount,
		Currency:       settlement.Currency,
	}

	task, err := queue.NewSettlementProcessTask(payload)
	if err != nil {
		return err
	}

	if _, err := s.queueClient.EnqueueContext(ctx, task); err != nil {
		log.Error().Err(err).Msg("failed to enqueue settlement task")
		return err
	}

	log.Info().Str("id", settlement.ID).Msg("settlement enqueued for processing")
	return nil
}