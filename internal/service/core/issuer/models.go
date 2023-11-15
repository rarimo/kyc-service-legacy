package issuer

import (
	"github.com/pkg/errors"
	"github.com/rarimo/issuer/resources"
)

var (
	ErrUnexpectedStatusCode = errors.New("unexpected status code")
)

const (
	issueEndpoint = "/private/claims/issue/{identity_id}/{claim_type}"
)

const (
	identityIDPathParam = "identity_id"
	claimTypePathParam  = "claim_type"
)

type ClaimType string

func (c ClaimType) String() string {
	return string(c)
}

const (
	ClaimTypeNaturalPerson     ClaimType = "NaturalPerson"
	ClaimTypeIdentityProviders ClaimType = "IdentityProviders"

	EmptyStringField  = "none"
	EmptyIntegerField = 0
)

type IdentityProviderName string

func (ipn IdentityProviderName) String() string {
	return string(ipn)
}

const (
	UnstoppableDomainsProviderName IdentityProviderName = "UnstoppableDomains"
	CivicProviderName              IdentityProviderName = "Civic"
	GitcoinProviderName            IdentityProviderName = "GitcoinPassport"
	WorldCoinProviderName          IdentityProviderName = "Worldcoin"
)

type IsNaturalPersonCredentialSubject struct {
	IsNatural string `json:"is_natural"`
}

type IdentityProvidersCredentialSubject struct {
	Provider                 IdentityProviderName `json:"provider"`
	IsNatural                int64                `json:"is_natural"`
	Address                  string               `json:"address"`
	GitcoinPassportScore     string               `json:"gitcoin_passport_score"`
	WorldCoinScore           string               `json:"worldcoin_score"`
	UnstoppableDomain        string               `json:"unstoppable_domain"`
	CivicGatekeeperNetworkID int64                `json:"civic_gatekeeper_network_id"`
	KYCAdditionalData        string               `json:"kyc_additional_data"`
}

func NewEmptyIdentityProvidersCredentialSubject() *IdentityProvidersCredentialSubject {
	return &IdentityProvidersCredentialSubject{
		Provider:                 EmptyStringField,
		IsNatural:                EmptyIntegerField,
		Address:                  EmptyStringField,
		GitcoinPassportScore:     EmptyStringField,
		WorldCoinScore:           EmptyStringField,
		UnstoppableDomain:        EmptyStringField,
		CivicGatekeeperNetworkID: EmptyIntegerField,
		KYCAdditionalData:        EmptyStringField,
	}
}

type IssueClaimResponse struct {
	Data IssueClaimResponseData `json:"data"`
}

type IssueClaimResponseData struct {
	resources.Key
}
