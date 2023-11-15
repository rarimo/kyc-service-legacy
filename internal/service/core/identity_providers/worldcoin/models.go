package worldcoin

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/pkg/errors"

	"github.com/rarimo/kyc-service-legacy/resources"
)

const (
	defaultBaseURL   = "https://id.worldcoin.org"
	userInfoEndpoint = "/userinfo"

	likelyHumanStrong = "strong"
)

type (
	// VerificationData is data that is required by Worldcoin to verify a user
	VerificationData resources.WorldcoinData

	// UserInfo is data that is returned by Worldcoin by /userinfo endpoint.
	UserInfo struct {
		Sub          string       `json:"sub"`
		HumanityInfo HumanityInfo `json:"https://id.worldcoin.org/beta"`
		Email        string       `json:"email,omitempty"`
		Name         string       `json:"name,omitempty"`
		GivenName    string       `json:"given_name,omitempty"`
		FamilyName   string       `json:"family_name,omitempty"`
	}

	HumanityInfo struct {
		LikelyHuman    string `json:"likely_human"`
		CredentialType string `json:"credential_type"`
	}
)

var (
	ErrNotLikelyHuman = errors.New("not likely human")
)

func unexpectedStatusCode(statusCode int) error {
	return errors.New(fmt.Sprintf("received unexpected status code: %d", statusCode))
}

func (v VerificationData) Validate() error {
	return validation.Errors{
		"id_token": validation.Validate(v.IdToken, validation.Required),
	}.Filter()
}
