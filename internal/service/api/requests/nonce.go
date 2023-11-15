package requests

import (
	"encoding/json"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pkg/errors"

	"github.com/rarimo/kyc-service-legacy/resources"
)

type NonceRequest struct {
	Address common.Address
}

type nonceRequest struct {
	Data resources.Nonce `json:"data"`
}

func NewNonceRequest(r *http.Request) (*NonceRequest, error) {
	var requestBody nonceRequest

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}

	if err := requestBody.validate(); err != nil {
		return nil, err
	}

	return requestBody.parse(), nil
}

func (r *nonceRequest) validate() error {
	return validation.Errors{
		"/data/type": validation.Validate(
			r.Data.Type, validation.Required, validation.In(resources.NONCE_REQUEST),
		),
		"/data/attributes/address": validation.Validate(
			r.Data.Attributes.Address, validation.Required, validation.By(MustBeEthAddress),
		),
	}.Filter()
}

func (r *nonceRequest) parse() *NonceRequest {
	return &NonceRequest{
		Address: common.HexToAddress(r.Data.Attributes.Address),
	}
}

func MustBeEthAddress(src interface{}) error {
	raw, ok := src.(string)
	if !ok {
		return errors.New("it is not a string")
	}

	if !common.IsHexAddress(raw) {
		return errors.New("it is not a valid ethereum address")
	}

	return nil
}
