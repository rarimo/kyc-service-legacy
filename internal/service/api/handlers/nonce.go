package handlers

import (
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"

	"github.com/rarimo/kyc-service-legacy/internal/service/api/requests"
	"github.com/rarimo/kyc-service-legacy/internal/service/api/responses"
)

func GetNonce(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewNonceRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	nonce, err := KYCService(r).NewNonce(req)
	switch {
	case err != nil:
		Log(r).WithError(err).Error("Failed to create new nonce")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, responses.NewNonce(nonce.Nonce))
}
