package handlers

import (
	"errors"
	"github.com/rarimo/kyc-service-legacy/internal/service/core"
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"

	"github.com/rarimo/kyc-service-legacy/internal/service/api/requests"
	"github.com/rarimo/kyc-service-legacy/internal/service/api/responses"
)

func GetVerifyStatus(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewVerifyStatusRequest(r)
	if err != nil {
		Log(r).WithError(err).Error("failed to parse request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	user, err := KYCService(r).GetVerifyStatus(req)
	switch {
	case errors.Is(err, core.ErrUserNotFound):
		ape.RenderErr(w, problems.NotFound())
		return
	case err != nil:
		Log(r).WithError(err).Error("failed to get user")
		ape.RenderErr(w, problems.InternalError())
		return
	}

	ape.Render(w, responses.NewVerifyStatus(user))
}
