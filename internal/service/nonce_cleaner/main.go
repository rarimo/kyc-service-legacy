package nonceCleaner

import (
	"context"
	"time"

	"github.com/rarimo/kyc-service-legacy/internal/config"
	"github.com/rarimo/kyc-service-legacy/internal/data"
	"github.com/rarimo/kyc-service-legacy/internal/data/pg"

	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/running"
)

func NewNonceCleaner(cfg config.Config) *NonceCleaner {
	return &NonceCleaner{
		storage: pg.NewMasterQ(cfg.DB()).NonceQ(),
		logger:  cfg.Log(),
	}
}

type NonceCleaner struct {
	storage data.NonceQ
	logger  *logan.Entry
}

func (s *NonceCleaner) Run(ctx context.Context) {
	running.WithBackOff(ctx,
		s.logger,
		"nonce-cleaner",
		func(ctx context.Context) error {
			return s.storage.WhereExpiresAtLt(time.Now()).Delete()
		},
		30*time.Minute,
		1*time.Second,
		5*time.Second,
	)
}

func Run(ctx context.Context, cfg config.Config) {
	NewNonceCleaner(cfg).Run(ctx)
}
