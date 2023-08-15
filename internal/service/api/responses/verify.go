package responses

import (
	"fmt"
	"github.com/google/jsonapi"
	"github.com/google/uuid"
	"gitlab.com/rarimo/identity/kyc-service/internal/service/core"
	providers "gitlab.com/rarimo/identity/kyc-service/internal/service/core/identity_providers"
	"gitlab.com/rarimo/identity/kyc-service/internal/service/core/identity_providers/civic"
	gcpsp "gitlab.com/rarimo/identity/kyc-service/internal/service/core/identity_providers/gitcoin_passport"
	"gitlab.com/rarimo/identity/kyc-service/internal/service/core/identity_providers/worldcoin"
	"net/http"

	"gitlab.com/rarimo/identity/kyc-service/resources"
)

const (
	InvalidAccessToken = iota + 1
	InvalidUsersSignature
	NonceNotFound
	NotLikelyHuman
	ScoreIsTooLow
	InvalidGatewayToken
)

const (
	EthAddressTaken = iota + 1
	IdentityVerified
	DuplicateProviderData
)

var unauthorizedErrorCodes = map[string]int{
	providers.ErrInvalidAccessToken.Error():    InvalidAccessToken,
	providers.ErrInvalidUsersSignature.Error(): InvalidUsersSignature,
	providers.ErrNonceNotFound.Error():         NonceNotFound,
	worldcoin.ErrNotLikelyHuman.Error():        NotLikelyHuman,
	gcpsp.ErrScoreIsTooLow.Error():             ScoreIsTooLow,
	civic.ErrInvalidGatewayToken.Error():       InvalidGatewayToken,
}

var conflictErrorCodes = map[string]int{
	core.ErrUserAlreadyVerifiedByEthAddress.Error(): EthAddressTaken,
	core.ErrUserAlreadyVerifiedByIdentityID.Error(): IdentityVerified,
	core.ErrDuplicatedProviderData.Error():          DuplicateProviderData,
}

type VerifyIDResponse struct {
	Data VerifyID `json:"data"`
}

type VerifyID struct {
	resources.Key
}

func NewVerify(id uuid.UUID) VerifyIDResponse {
	return VerifyIDResponse{
		Data: VerifyID{
			Key: resources.Key{
				ID:   id.String(),
				Type: resources.VERIFICATION_ID,
			},
		},
	}
}

func Unauthorized(cause string) *jsonapi.ErrorObject {
	return &jsonapi.ErrorObject{
		Title:  http.StatusText(http.StatusUnauthorized),
		Status: fmt.Sprintf("%d", http.StatusUnauthorized),
		Code:   fmt.Sprintf("%d", unauthorizedErrorCodes[cause]),
		Detail: cause,
	}
}

func Conflict(cause string) *jsonapi.ErrorObject {
	return &jsonapi.ErrorObject{
		Title:  http.StatusText(http.StatusConflict),
		Status: fmt.Sprintf("%d", http.StatusConflict),
		Code:   fmt.Sprintf("%d", conflictErrorCodes[cause]),
		Detail: cause,
	}
}
