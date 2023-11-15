package identityproviders

import (
	"github.com/rarimo/kyc-service-legacy/internal/data"
	"github.com/rarimo/kyc-service-legacy/internal/service/core/issuer"
)

type IdentityProvider interface {
	Verify(user *data.User, verifyProviderDataRaw []byte) (*issuer.IdentityProvidersCredentialSubject, []byte, error)
}
