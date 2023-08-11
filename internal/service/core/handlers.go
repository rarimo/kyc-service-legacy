package core

import (
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"gitlab.com/rarimo/identity/kyc-service/internal/crypto"
	"gitlab.com/rarimo/identity/kyc-service/internal/data"
	"gitlab.com/rarimo/identity/kyc-service/internal/service/api/requests"
	"gitlab.com/rarimo/identity/kyc-service/internal/service/core/issuer"
)

func (k *kycService) NewVerifyRequest(req *requests.VerifyRequest) (*data.User, error) {
	prevUser, err := k.db.UsersQ().WhereIdentityID(data.NewIdentityID(req.IdentityID)).Get()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user from db with the same identityID")
	}
	if prevUser != nil && prevUser.Status != data.UserStatusUnverified {
		return nil, ErrUserAlreadyVerifiedByIdentityID
	}

	newUser := data.User{
		ID:         uuid.New(),
		Status:     data.UserStatusInitialized,
		CreatedAt:  time.Now(),
		IdentityID: data.NewIdentityID(req.IdentityID),
	}

	credentialSubject, err := k.verifyUser(req, &newUser)
	if err != nil {
		return nil, errors.Wrap(err, "failed to verify user")
	}

	err = k.db.Transaction(func() error {
		if err = k.db.UsersQ().Upsert(&newUser); err != nil {
			return errors.Wrap(err, "failed to upsert user into db")
		}

		// "1" == true
		credentialSubject.IsNatural = 1
		if newUser.Status == data.UserStatusVerified {
			_, err := k.issuer.IssueClaim(
				newUser.IdentityID.ID,
				issuer.ClaimTypeIdentityProviders,
				credentialSubject,
			)
			if err != nil {
				return errors.Wrap(err, "failed to issue claim")
			}
		}

		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute db transaction")
	}

	return &newUser, nil
}

func (k *kycService) NewNonce(req *requests.NonceRequest) (*data.Nonce, error) {
	if err := k.db.NonceQ().WhereEthAddress(req.Address).Delete(); err != nil {
		return nil, errors.Wrap(err, "failed to delete old nonce")
	}

	newNonce, err := crypto.NewNonce()
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate nonce")
	}

	nonce := &data.Nonce{
		EthAddress: req.Address,
		Nonce:      newNonce,
		ExpiresAt:  time.Now().Add(k.nonceLifetime),
	}

	err = k.db.NonceQ().Insert(nonce)
	if err != nil {
		return nil, errors.Wrap(err, "failed to insert nonce into db")
	}

	return nonce, nil
}

func (k *kycService) GetVerifyStatus(req *requests.VerifyStatusRequest) (*data.User, error) {
	user, err := k.db.UsersQ().WhereID(req.VerifyID).Get()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user from db")
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return user, nil
}

func (k *kycService) verifyUser(
	req *requests.VerifyRequest, newUser *data.User,
) (*issuer.IdentityProvidersCredentialSubject, error) {
	credSubject, providerHash, err := k.identityProviders[req.IdentityProviderName].Verify(newUser, req.ProviderData)
	if err != nil {
		return nil, errors.Wrap(err, "failed to verify with identity provider")
	}

	user, err := k.db.UsersQ().WhereProviderHash(providerHash).Get()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user from db with the same providerDataHash")
	}
	if user != nil {
		return nil, ErrDuplicatedProviderData
	}

	newUser.ProviderHash = providerHash

	return credSubject, nil
}
