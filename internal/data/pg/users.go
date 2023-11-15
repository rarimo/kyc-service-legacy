package pg

import (
	"database/sql"

	"github.com/fatih/structs"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/kit/pgdb"

	sq "github.com/Masterminds/squirrel"
	"github.com/ethereum/go-ethereum/common"

	"github.com/rarimo/kyc-service-legacy/internal/data"
)

type usersQ struct {
	db  *pgdb.DB
	sel sq.SelectBuilder
}

func NewUsersQ(db *pgdb.DB) data.UsersQ {
	return &usersQ{
		db:  db,
		sel: sq.Select("*").From(usersTableName),
	}
}

func (q *usersQ) New() data.UsersQ {
	return NewUsersQ(q.db.Clone())
}

func (q *usersQ) Select() ([]data.User, error) {
	var result []data.User

	err := q.db.Select(&result, q.sel)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to select rows")
	}

	return result, nil
}

func (q *usersQ) Get() (*data.User, error) {
	var result data.User

	err := q.db.Get(&result, q.sel)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to select rows")
	}

	return &result, nil
}

func (q *usersQ) Insert(user *data.User) error {
	clauses := structs.Map(user)
	clauses[identityIDColumnName] = user.IdentityID

	err := q.db.Exec(
		sq.Insert(usersTableName).
			SetMap(clauses),
	)
	if err != nil {
		return errors.Wrap(err, "failed to insert rows")
	}

	return nil
}

func (q *usersQ) Upsert(user *data.User) error {
	clauses := structs.Map(user)
	clauses[identityIDColumnName] = user.IdentityID

	err := q.db.Exec(
		sq.Insert(usersTableName).
			SetMap(clauses).
			Suffix("ON CONFLICT (identity_id) " +
				"DO UPDATE SET " +
				"status = EXCLUDED.status, " +
				"eth_address = EXCLUDED.eth_address, " +
				"provider_data = EXCLUDED.provider_data, " +
				"created_at = EXCLUDED.created_at, " +
				"provider_hash = EXCLUDED.provider_hash"),
	)
	if err != nil {
		return errors.Wrap(err, "failed to upsert rows")
	}

	return nil
}

func (q *usersQ) Update(user *data.User) error {
	err := q.db.Exec(
		sq.Update(usersTableName).
			SetMap(structs.Map(user)).
			Where(sq.Eq{idColumnName: user.ID}),
	)
	if err != nil {
		return errors.Wrap(err, "failed to update rows")
	}

	return nil
}

func (q *usersQ) Sort(sort pgdb.SortedOffsetPageParams) data.UsersQ {
	q.sel = sort.ApplyTo(q.sel, sortColumns)
	return q
}

func (q *usersQ) WhereID(id uuid.UUID) data.UsersQ {
	q.sel = q.sel.Where(sq.Eq{idColumnName: id})
	return q
}

func (q *usersQ) WhereStatus(status data.UserStatus) data.UsersQ {
	q.sel = q.sel.Where(sq.Eq{statusColumnName: status})
	return q
}

func (q *usersQ) WhereEthAddress(address common.Address) data.UsersQ {
	q.sel = q.sel.Where(sq.Eq{ethAddressColumnName: address})
	return q
}

func (q *usersQ) WhereIdentityID(identityID *data.IdentityID) data.UsersQ {
	q.sel = q.sel.Where(sq.Eq{identityIDColumnName: identityID})
	return q
}

func (q *usersQ) WhereProviderHash(providerHash []byte) data.UsersQ {
	q.sel = q.sel.Where(sq.Eq{providerHashColumnName: providerHash})
	return q
}
