package requests

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	iden3core "github.com/iden3/go-iden3-core"
	"github.com/pkg/errors"

	providers "github.com/rarimo/kyc-service-legacy/internal/service/core/identity_providers"
	"github.com/rarimo/kyc-service-legacy/resources"
)

type VerifyRequest struct {
	IdentityProviderName providers.IdentityProviderName
	ProviderData         []byte
	IdentityID           iden3core.ID
}

type verifyRequest struct {
	IdentityProviderName string
	resources.VerifyRequest
}

func NewVerify(r *http.Request) (*VerifyRequest, error) {
	requestBody := verifyRequest{
		IdentityProviderName: chi.URLParam(r, IdentityProviderPathParam),
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return nil, errors.Wrap(err, "failed to decode json request body")
	}

	if err := requestBody.validate(); err != nil {
		return nil, err
	}

	return requestBody.parse(), nil
}

func (req *verifyRequest) validate() error {
	return validation.Errors{
		fmt.Sprintf("path/%s", IdentityProviderPathParam): validation.Validate(
			req.IdentityProviderName, validation.Required, validation.By(MustBeValidIdentityProviderName),
		),
		"body/data/type": validation.Validate(
			req.Data.Type, validation.Required, validation.In(resources.VERIFY),
		),
		"body/data/attributes/provider_data": validation.Validate(
			req.Data.Attributes.ProviderData, validation.Required, validation.By(MustBeValidJSONOrEmpty),
		),
		"body/data/attributes/identity_id": validation.Validate(
			req.Data.Attributes.IdentityId, validation.Required, validation.By(MustBeValidID),
		),
	}.Filter()
}

func (req *verifyRequest) parse() *VerifyRequest {
	identityID, _ := iden3core.IDFromString(req.Data.Attributes.IdentityId)

	return &VerifyRequest{
		IdentityProviderName: providers.IdentityProviderNames[req.IdentityProviderName],
		ProviderData:         req.Data.Attributes.ProviderData,
		IdentityID:           identityID,
	}
}

func MustBeValidIdentityProviderName(src interface{}) error {
	nameRaw, ok := src.(string)
	if !ok {
		return errors.New("it is not a string")
	}

	if _, ok := providers.IdentityProviderNames[nameRaw]; !ok {
		return errors.New("it is not a valid identity provider name")
	}

	return nil
}

func MustBeValidJSONOrEmpty(src interface{}) error {
	dataRaw, ok := src.(json.RawMessage)
	if !ok {
		return errors.New("it is not a byte array")
	}

	if len(dataRaw) == 0 {
		return nil
	}

	var data interface{}
	if err := json.Unmarshal(dataRaw, &data); err != nil {
		return errors.Wrap(err, "it is not a valid json")
	}

	return nil
}

func MustBeValidID(src interface{}) error {
	identifierRawBase58, ok := src.(string)
	if !ok {
		return errors.New("it is not a string")
	}

	id, err := iden3core.IDFromString(identifierRawBase58)
	if err != nil {
		return errors.New("it is not an iden3 identity id")
	}

	_, err = iden3core.ParseDIDFromID(id)
	if err != nil {
		return errors.New("can't parse to did")
	}

	return nil
}
