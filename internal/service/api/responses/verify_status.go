package responses

import (
	"github.com/rarimo/kyc-service-legacy/internal/data"
	"github.com/rarimo/kyc-service-legacy/resources"
)

func NewVerifyStatus(user *data.User) resources.VerifyStatusRequest {
	return resources.VerifyStatusRequest{
		Data: resources.VerifyStatus{
			Key: resources.Key{
				ID:   user.ID.String(),
				Type: resources.VERIFICATION_ID,
			},
			Attributes: resources.VerifyStatusAttributes{
				Status: user.Status.String(),
			},
		},
	}
}
