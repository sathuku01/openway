package queue

import (
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

func NewRedisClient(opt asynq.RedisConnOpt) *asynq.Client {
	client := asynq.NewClient(opt)
	log.Info().Msg("redis queue client initialized")
	return client
}

func NewBackgroundServer(opt asynq.RedisConnOpt) *asynq.Server {
	srv := asynq.NewServer(
		opt,
		asynq.Config{
			Concurrency: 10,
			Queues: map[string]int{
				"critical": 6,
				"default":  3,
				"low":      1,
			},
		},
	)
	log.Info().Msg("redis background worker server initialized")
	return srv
}