package pg

import (
	"gitlab.com/distributed_lab/kit/pgdb"

	"github.com/rarimo/kyc-service-legacy/internal/data"
)

const (
	usersTableName = "users"
)

const (
	idColumnName           = "id"
	statusColumnName       = "status"
	ethAddressColumnName   = "eth_address"
	identityIDColumnName   = "identity_id"
	providerHashColumnName = "provider_hash"
	expiresAtColumnName    = "expires_at"
)

var sortColumns = map[string]string{
	statusColumnName: statusColumnName,
}

type masterQ struct {
	db *pgdb.DB
}

func NewMasterQ(db *pgdb.DB) data.MasterQ {
	return &masterQ{
		db: db,
	}
}

func (q *masterQ) New() data.MasterQ {
	return NewMasterQ(q.db.Clone())
}

func (q *masterQ) UsersQ() data.UsersQ {
	return NewUsersQ(q.db)
}

func (q *masterQ) NonceQ() data.NonceQ {
	return NewNonceQ(q.db)
}

func (q *masterQ) Transaction(fn func() error) error {
	return q.db.Transaction(func() error {
		return fn()
	})
}
