package gcpsp

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"gitlab.com/distributed_lab/logan/v3/errors"

	"github.com/rarimo/kyc-service-legacy/internal/service/api/requests"
	"github.com/rarimo/kyc-service-legacy/resources"
)

const (
	defaultBaseURL         = "https://api.scorer.gitcoin.co/registry"
	submitPassportEndpoint = "/submit-passport"
	scoreEndpoint          = "/score"
	defaultRetryCount      = 3
)

const (
	statusDone       = "DONE"
	statusProcessing = "PROCESSING"
)

const (
	// retryInterval is a time interval between request retries
	retryInterval = time.Second * 3
)

type (
	// VerificationData is a data that is required by Gitcoin Passport to verify a user
	VerificationData resources.GitcoinPassportData

	// SubmitPassportResponse is a data that is returned by Gitcoin Passport with status of request
	SubmitPassportResponse struct {
		Address string `json:"address"`
		Score   string `json:"score"`
		Status  string `json:"status"`
	}

	// GetStampsResponse is a data that is returned by Gitcoin Passport with stamps
	GetStampsResponse struct {
		Items []Stamp `json:"items"`
	}

	// SubmitPassportRequest is a data that is sent to Gitcoin Passport
	SubmitPassportRequest struct {
		ScorerId int    `json:"scorer_id"`
		Address  string `json:"address"`
	}

	// ProviderData is a data that is returned by Gitcoin Passport
	ProviderData struct {
		Address string  `json:"address"`
		Score   string  `json:"score"`
		Stamps  []Stamp `json:"stamps"`
	}

	// Stamp is a user data that is returned by Gitcoin Passport
	Stamp struct {
		Version    string `json:"version"`
		Credential struct {
			Type  []string `json:"type"`
			Proof struct {
				Jws                string `json:"jws"`
				Type               string `json:"type"`
				Created            string `json:"created"`
				ProofPurpose       string `json:"proofPurpose"`
				VerificationMethod string `json:"verificationMethod"`
			} `json:"proof"`
			Issuer            string   `json:"issuer"`
			Context           []string `json:"@context"`
			IssuanceDate      string   `json:"issuanceDate"`
			ExpirationDate    string   `json:"expirationDate"`
			CredentialSubject struct {
				Id      string `json:"id"`
				Hash    string `json:"hash"`
				Context []struct {
					Hash     string `json:"hash"`
					Provider string `json:"provider"`
				} `json:"@context"`
				Provider string `json:"provider"`
			} `json:"credentialSubject"`
		} `json:"credential"`
		Metadata *struct {
			Group    string `json:"group"`
			Platform struct {
				Id             string `json:"id"`
				Icon           string `json:"icon"`
				Name           string `json:"name"`
				Description    string `json:"description"`
				ConnectMessage string `json:"connectMessage"`
			} `json:"platform"`
			Name        string `json:"name"`
			Description string `json:"description"`
			Hash        string `json:"hash"`
		} `json:"metadata,omitempty"`
	}
)

var (
	ErrUnexpectedStatus     = errors.New("unexpected status")
	ErrUnexpectedStatusCode = errors.New("received unexpected status code")
	ErrScoreIsTooLow        = errors.New("score is too low")
)

// Validate is a method that validates VerificationData
func (v VerificationData) Validate() error {
	return validation.Errors{
		"signature": validation.Validate(v.Signature, validation.Required),
		"address":   validation.Validate(v.Address, validation.Required, validation.By(requests.MustBeEthAddress)),
	}.Filter()
}

func getStampsEndpoint(address string) string {
	return fmt.Sprintf("stamps/%s?include_metadata=true", address)
}
