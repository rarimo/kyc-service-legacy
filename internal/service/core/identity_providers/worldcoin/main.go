package worldcoin

import (
	"encoding/json"
	"net/http"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/imroc/req/v3"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/logan/v3"

	"gitlab.com/rarimo/identity/kyc-service/internal/config"
	"gitlab.com/rarimo/identity/kyc-service/internal/data"
	providers "gitlab.com/rarimo/identity/kyc-service/internal/service/core/identity_providers"
	"gitlab.com/rarimo/identity/kyc-service/internal/service/core/issuer"
)

// Worldcoin is a struct that implements the identityproviders.IdentityProvider interface
type Worldcoin struct {
	*req.Client
}

// NewIdentityProvider returns a new Worldcoin struct
func NewIdentityProvider(log *logan.Entry, settings *config.WorldcoinSettings) *Worldcoin {
	if settings.BaseURL == "" {
		log.Warn("Worldcoin base url is not set, using default one")
		settings.BaseURL = defaultBaseURL
	}

	return &Worldcoin{Client: req.C().SetBaseURL(settings.BaseURL).SetLogger(log)}
}

// Verify verifies the user's identity
func (w *Worldcoin) Verify(
	user *data.User, verifyDataRaw []byte,
) (*issuer.IdentityProvidersCredentialSubject, []byte, error) {
	var verifyData VerificationData
	if err := json.Unmarshal(verifyDataRaw, &verifyData); err != nil {
		return nil, nil, errors.Wrap(err, "failed to unmarshal verification data")
	}

	if err := verifyData.Validate(); err != nil {
		return nil, nil, providers.ErrInvalidVerificationData
	}

	userInfo, err := w.retrieveUserInfo(verifyData.IdToken)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to retrieve user info")
	}

	if userInfo.HumanityInfo.LikelyHuman != likelyHumanStrong {
		return nil, nil, ErrNotLikelyHuman
	}

	userInfoRaw, err := json.Marshal(userInfo)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to marshal user info")
	}

	user.Status = data.UserStatusVerified
	user.ProviderData = userInfoRaw
	// as we don't have the user's eth address, we set it to the zero address
	user.EthAddress = nil

	credentialSubject := issuer.NewEmptyIdentityProvidersCredentialSubject()
	credentialSubject.Provider = issuer.WorldCoinProviderName
	credentialSubject.WorldCoinScore = likelyHumanStrong
	credentialSubject.KYCAdditionalData = string(userInfoRaw)

	return credentialSubject, crypto.Keccak256(
		[]byte(userInfo.Sub),
		providers.WorldCoinIdentityProvider.Bytes(),
	), nil
}

// retrieveUserInfo retrieves the user's info from the Worldcoin API
func (w *Worldcoin) retrieveUserInfo(accessToken string) (*UserInfo, error) {
	var userInfo UserInfo

	response, err := w.R().
		SetBearerAuthToken(accessToken).
		SetSuccessResult(&userInfo).
		Get(userInfoEndpoint)
	if err != nil {
		return nil, errors.Wrap(err, "failed to retrieve user info")
	}

	if response.StatusCode >= http.StatusBadRequest {
		if response.StatusCode == http.StatusUnauthorized {
			return nil, providers.ErrInvalidAccessToken
		}

		return nil, unexpectedStatusCode(response.StatusCode)
	}

	return &userInfo, nil
}
