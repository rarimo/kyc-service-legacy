package civic

import (
	"encoding/json"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	cryptoPkg "github.com/ethereum/go-ethereum/crypto"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/logan/v3"

	"github.com/rarimo/kyc-service-legacy/internal/config"
	"github.com/rarimo/kyc-service-legacy/internal/crypto"
	"github.com/rarimo/kyc-service-legacy/internal/data"
	providers "github.com/rarimo/kyc-service-legacy/internal/service/core/identity_providers"
	"github.com/rarimo/kyc-service-legacy/internal/service/core/identity_providers/civic/contracts"
	"github.com/rarimo/kyc-service-legacy/internal/service/core/issuer"
)

type Civic struct {
	log                           *logan.Entry
	masterQ                       data.MasterQ
	CaptchaGatekeeperNetworkId    *big.Int
	UniquenessGatekeeperNetworkId *big.Int
	ContractsVerifiers            map[ChainName]*contracts.GatewayVerifier

	skipSigCheck bool
}

func NewIdentityProvider(log *logan.Entry, masterQ data.MasterQ, config *config.Civic) (*Civic, error) {
	contractsVerifiers, err := newContractsVerifiers(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create contracts verifiers")
	}

	return &Civic{
		log:                           log,
		masterQ:                       masterQ,
		CaptchaGatekeeperNetworkId:    config.CaptchaGatekeeperNetworkId,
		UniquenessGatekeeperNetworkId: config.UniquenessGatekeeperNetworkId,
		ContractsVerifiers:            contractsVerifiers,
		skipSigCheck:                  config.SkipSigCheck,
	}, nil
}

func (c *Civic) Verify(
	user *data.User, verifyDataRaw []byte,
) (*issuer.IdentityProvidersCredentialSubject, []byte, error) {
	var verifyData VerificationData
	if err := json.Unmarshal(verifyDataRaw, &verifyData); err != nil {
		return nil, nil, errors.Wrap(err, "failed to unmarshal verification data")
	}

	if err := verifyData.Validate(); err != nil {
		return nil, nil, errors.Wrap(providers.ErrInvalidVerificationData, err.Error())
	}

	if err := c.verifySignature(verifyData, verifyData.Address); err != nil {
		return nil, nil, errors.Wrap(err, "failed to verify signature")
	}

	token, err := c.verifyGatewayToken(chainNameFromString[verifyData.ChainName], verifyData.Address)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to verify gateway token")
	}

	providerDataRaw, err := json.Marshal(ProviderData{
		Address: verifyData.Address,
	})
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to marshal provider data")
	}

	user.EthAddress = &verifyData.Address
	user.Status = data.UserStatusVerified
	user.ProviderData = providerDataRaw

	credentialSubject := issuer.NewEmptyIdentityProvidersCredentialSubject()
	credentialSubject.Provider = issuer.CivicProviderName
	credentialSubject.Address = verifyData.Address.String()
	credentialSubject.CivicGatekeeperNetworkID = token.Int64()

	return credentialSubject, cryptoPkg.Keccak256(
		verifyData.Address.Bytes(),
		providers.CivicIdentityProvider.Bytes(),
	), nil
}

func (c *Civic) verifySignature(verifyData VerificationData, userAddress common.Address) error {
	if c.skipSigCheck {
		return nil
	}

	nonce, err := c.masterQ.NonceQ().
		WhereEthAddress(verifyData.Address).
		WhereExpiresAtGt(time.Now()).
		Get()
	if err != nil {
		return errors.Wrap(err, "failed to get nonce by address")
	}
	if nonce == nil {
		return providers.ErrNonceNotFound
	}

	valid, err := crypto.VerifyEIP191Signature(
		verifyData.Signature,
		crypto.NonceToSignMessage(nonce.Nonce),
		userAddress,
	)
	if err != nil || !valid {
		return providers.ErrInvalidUsersSignature
	}

	if err = c.masterQ.NonceQ().WhereEthAddress(verifyData.Address).Delete(); err != nil {
		return errors.Wrap(err, "failed to delete nonce")
	}

	return nil
}

func (c *Civic) verifyGatewayToken(chainName ChainName, userAddress common.Address) (*big.Int, error) {
	verifier, ok := c.ContractsVerifiers[chainName]
	if !ok {
		return nil, ErrVerifierNotFound
	}

	validCaptcha, err := verifier.VerifyToken(nil, userAddress, c.CaptchaGatekeeperNetworkId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to verify gateway token on chain")
	}
	if validCaptcha {
		return c.CaptchaGatekeeperNetworkId, nil
	}

	validUniqueness, err := verifier.VerifyToken(nil, userAddress, c.UniquenessGatekeeperNetworkId)
	if err != nil {
		return nil, errors.Wrap(err, "failed to verify gateway token on chain")
	}
	if validUniqueness {
		return c.UniquenessGatekeeperNetworkId, nil
	}

	return nil, ErrInvalidGatewayToken
}

func newContractsVerifiers(config *config.Civic) (map[ChainName]*contracts.GatewayVerifier, error) {
	ethereumVerifier, err := contracts.NewGatewayVerifier(config.GatewayTokenContract, config.EthereumRpc)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Ethereum gateway verifier")
	}

	polygonVerifier, err := contracts.NewGatewayVerifier(config.GatewayTokenContract, config.PolygonRpc)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Polygon gateway verifier")
	}

	arbitrumVerifier, err := contracts.NewGatewayVerifier(config.GatewayTokenContract, config.ArbitrumRpc)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create Arbitrum gateway verifier")
	}

	xdcVerifier, err := contracts.NewGatewayVerifier(config.GatewayTokenContract, config.XDCRpc)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create XDC gateway verifier")
	}

	return map[ChainName]*contracts.GatewayVerifier{
		EthereumChainName: ethereumVerifier,
		PolygonChainName:  polygonVerifier,
		ArbitrumChainName: arbitrumVerifier,
		XDCChainName:      xdcVerifier,
	}, nil
}
