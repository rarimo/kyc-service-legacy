package api

import (
	"context"
	"net"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/kit/copus/types"
	"gitlab.com/distributed_lab/logan/v3"
	"gitlab.com/distributed_lab/logan/v3/errors"

	"github.com/rarimo/kyc-service-legacy/internal/config"
	"github.com/rarimo/kyc-service-legacy/internal/service/core"
)

type service struct {
	log        *logan.Entry
	copus      types.Copus
	listener   net.Listener
	kycService core.KYCService
}

func newService(cfg config.Config, ctx context.Context) *service {
	kycService, err := core.NewKYCService(cfg, ctx)
	if err != nil {
		panic(errors.Wrap(err, "failed to create new KYC service"))
	}

	return &service{
		log:        cfg.Log(),
		copus:      cfg.Copus(),
		listener:   cfg.Listener(),
		kycService: kycService,
	}
}

func Run(ctx context.Context, cfg config.Config) {
	svc := newService(cfg, ctx)

	if err := svc.copus.RegisterChi(svc.router()); err != nil {
		panic(errors.Wrap(err, "cop failed"))
	}

	svc.log.Info("Service started")
	ape.Serve(ctx, svc.router(), cfg, ape.ServeOpts{})
}
