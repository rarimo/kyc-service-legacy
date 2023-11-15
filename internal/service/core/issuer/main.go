package issuer

import (
	"encoding/json"

	iden3core "github.com/iden3/go-iden3-core"
	"github.com/imroc/req/v3"
	"github.com/pkg/errors"
	"github.com/rarimo/issuer/resources"
	"gitlab.com/distributed_lab/logan/v3"

	"github.com/rarimo/kyc-service-legacy/internal/config"
)

type Issuer interface {
	IssueClaim(
		identityID iden3core.ID,
		claimType ClaimType,
		credentialSubject interface{},
	) (*IssueClaimResponse, error)
}

type issuer struct {
	*req.Client
}

func New(log *logan.Entry, config *config.Issuer) Issuer {
	return &issuer{
		Client: req.C().
			SetBaseURL(config.BaseURL).
			SetLogger(log),
	}
}

func (is *issuer) IssueClaim(
	identityID iden3core.ID,
	claimType ClaimType,
	credentialsSubject interface{},
) (*IssueClaimResponse, error) {
	var result IssueClaimResponse

	requestBody, err := compactIssueClaimBody(credentialsSubject)
	if err != nil {
		return nil, errors.Wrap(err, "failed to compact issue claim request body")
	}

	response, err := is.R().
		SetBodyJsonMarshal(requestBody).
		SetSuccessResult(result).
		SetPathParams(map[string]string{
			identityIDPathParam: identityID.String(),
			claimTypePathParam:  claimType.String(),
		}).
		Post(issueEndpoint)
	if err != nil {
		return nil, errors.Wrap(err, "failed to send post request")
	}

	if response.StatusCode >= 299 {
		return nil, errors.Wrapf(ErrUnexpectedStatusCode, response.String())
	}

	return &result, nil
}

func compactIssueClaimBody(credentialSubject interface{}) (interface{}, error) {
	credentialsSubjectRaw, err := json.Marshal(credentialSubject)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal credentials subject")
	}

	return &resources.IssueClaimRequest{
		Data: resources.IssueClaim{
			Key: resources.Key{
				Type: resources.CLAIM_ISSUE,
			},
			Attributes: resources.IssueClaimAttributes{
				CredentialSubject: credentialsSubjectRaw,
			},
		},
	}, nil
}
