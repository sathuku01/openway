package queue

import (
	"context"
	"encoding/json"

	"openway/internal/repository"
	"openway/internal/web3"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type SettlementWorker struct {
	repo   *repository.SettlementRepository
	signer *web3.Signer
}

func NewSettlementWorker(repo *repository.SettlementRepository, signer *web3.Signer) *SettlementWorker {
	return &SettlementWorker{repo: repo, signer: signer}
}

func (w *SettlementWorker) RegisterHandlers() *asynq.ServeMux {
	mux := asynq.NewServeMux()
	mux.HandleFunc(TypeSettlementProcess, w.handleSettlementProcess)
	return mux
}

func (w *SettlementWorker) handleSettlementProcess(ctx context.Context, t *asynq.Task) error {
	var p SettlementProcessPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return err
	}

	log.Info().Str("id", p.SettlementID).Msg("processing settlement")
	
	var dummyIntentHash [32]byte
	var telcoRefHash [32]byte
	copy(telcoRefHash[:], []byte(p.TelcoReference))

	_, err := w.signer.SignReleaseMessage(dummyIntentHash, telcoRefHash)
	if err != nil {
		log.Error().Err(err).Msg("failed to sign release message")
		return err
	}

	err = w.repo.UpdateStatus(ctx, p.SettlementID, repository.StatusProcessed, dummyIntentHash[:])
	if err != nil {
		log.Error().Err(err).Msg("failed to update settlement status")
		return err
	}

	log.Info().Str("id", p.SettlementID).Msg("settlement processed and signed")
	return nil
}