package handlers

import (
	"net/http"

	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"

	providers "gitlab.com/rarimo/identity/kyc-service/internal/service/core/identity_providers"
	"gitlab.com/rarimo/identity/kyc-service/internal/service/core/identity_providers/civic"
	gcpsp "gitlab.com/rarimo/identity/kyc-service/internal/service/core/identity_providers/gitcoin_passport"
	"gitlab.com/rarimo/identity/kyc-service/internal/service/core/identity_providers/worldcoin"

	"gitlab.com/rarimo/identity/kyc-service/internal/data"
	"gitlab.com/rarimo/identity/kyc-service/internal/service/api/requests"
	"gitlab.com/rarimo/identity/kyc-service/internal/service/api/responses"
	"gitlab.com/rarimo/identity/kyc-service/internal/service/core"
)

func Verify(w http.ResponseWriter, r *http.Request) {
	req, err := requests.NewVerify(r)
	if err != nil {
		Log(r).WithField("reason", err).Debug("Bad request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	user, err := KYCService(r).NewVerifyRequest(req)
	switch {
	case errors.Is(err, providers.ErrInvalidVerificationData):
		Log(r).WithField("reason", err).
			WithField("identity-provider", req.IdentityProviderName).
			WithField("provider-data", string(req.ProviderData)).
			Debug("Bad request")
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	case errors.Is(err, core.ErrUserAlreadyVerifiedByEthAddress),
		errors.Is(err, core.ErrUserAlreadyVerifiedByIdentityID),
		errors.Is(err, core.ErrDuplicatedProviderData):
		Log(r).WithField("reason", err).Debug("Conflict")
		ape.RenderErr(w, responses.Conflict(errors.Cause(err).Error()))
		return
	case isUnauthorizedError(err):
		Log(r).WithField("reason", err).
			WithField("identity-provider", req.IdentityProviderName).
			WithField("provider-data", string(req.ProviderData)).
			Debug("Unauthorized")
		ape.RenderErr(w, responses.Unauthorized(errors.Cause(err).Error()))
		return
	case err != nil:
		Log(r).WithError(err).
			WithField("identity-provider", req.IdentityProviderName).
			WithField("provider-data", string(req.ProviderData)).
			Error("Failed to init new verification")
		ape.RenderErr(w, problems.InternalError())
		return
	case user.Status == data.UserStatusVerified:
		w.WriteHeader(http.StatusNoContent)
		return
	}

	ape.Render(w, responses.NewVerify(user.ID))
}

func isUnauthorizedError(err error) bool {
	return errors.Is(err, providers.ErrInvalidAccessToken) ||
		errors.Is(err, providers.ErrInvalidUsersSignature) ||
		errors.Is(err, providers.ErrNonceNotFound) ||
		errors.Is(err, worldcoin.ErrNotLikelyHuman) ||
		errors.Is(err, gcpsp.ErrScoreIsTooLow) ||
		errors.Is(err, civic.ErrInvalidGatewayToken)
}
