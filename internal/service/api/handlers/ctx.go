package handlers

import (
	"context"
	"net/http"

	"gitlab.com/distributed_lab/logan/v3"

	"github.com/rarimo/kyc-service-legacy/internal/service/core"
)

type ctxKey int

const (
	logCtxKey ctxKey = iota
	kycServiceCtxKey
)

func CtxLog(entry *logan.Entry) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, logCtxKey, entry)
	}
}

func Log(r *http.Request) *logan.Entry {
	return r.Context().Value(logCtxKey).(*logan.Entry)
}

func CtxKYCService(kycService core.KYCService) func(context.Context) context.Context {
	return func(ctx context.Context) context.Context {
		return context.WithValue(ctx, kycServiceCtxKey, kycService)
	}
}

func KYCService(r *http.Request) core.KYCService {
	return r.Context().Value(kycServiceCtxKey).(core.KYCService).New()
}
