package pg

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/rarimo/kyc-service-legacy/internal/data"

	sq "github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	"github.com/pkg/errors"
	"gitlab.com/distributed_lab/kit/pgdb"
)

const nonceTableName = "nonce"

type nonceQ struct {
	db  *pgdb.DB
	sel sq.SelectBuilder
	del sq.DeleteBuilder
}

func NewNonceQ(db *pgdb.DB) data.NonceQ {
	return &nonceQ{
		db:  db,
		sel: sq.Select("*").From(nonceTableName),
		del: sq.Delete(nonceTableName),
	}
}

func (q *nonceQ) New() data.NonceQ {
	return NewNonceQ(q.db.Clone())
}

func (q *nonceQ) Get() (*data.Nonce, error) {
	var result data.Nonce

	err := q.db.Get(&result, q.sel)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to get rows")
	}

	return &result, nil
}

func (q *nonceQ) Select() ([]data.Nonce, error) {
	var result []data.Nonce

	err := q.db.Select(&result, q.sel)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "failed to select rows")
	}

	return result, nil
}

func (q *nonceQ) Insert(value *data.Nonce) error {
	err := q.db.Get(&value.ID,
		sq.Insert(nonceTableName).
			SetMap(structs.Map(value)).
			Suffix(fmt.Sprintf("returning %s", idColumnName)),
	)
	if err != nil {
		return errors.Wrap(err, "failed to insert rows")
	}

	return nil
}

func (q *nonceQ) Update(value *data.Nonce) error {
	err := q.db.Exec(
		sq.Update(nonceTableName).
			SetMap(structs.Map(value)).
			Where(sq.Eq{idColumnName: value.ID}),
	)
	if err != nil {
		return errors.Wrap(err, "failed to update rows")
	}

	return nil
}

func (q *nonceQ) Delete() error {
	err := q.db.Exec(q.del)
	if err != nil {
		return errors.Wrap(err, "failed to delete rows")
	}

	return nil
}

func (q *nonceQ) WhereEthAddress(addresses ...common.Address) data.NonceQ {
	q.sel = q.sel.Where(sq.Eq{ethAddressColumnName: addresses})
	q.del = q.del.Where(sq.Eq{ethAddressColumnName: addresses})
	return q
}

func (q *nonceQ) WhereExpiresAtLt(expiresAt time.Time) data.NonceQ {
	q.sel = q.sel.Where(sq.Lt{expiresAtColumnName: expiresAt.UTC()})
	q.del = q.del.Where(sq.Lt{expiresAtColumnName: expiresAt.UTC()})
	return q
}

func (q *nonceQ) WhereExpiresAtGt(expiresAt time.Time) data.NonceQ {
	q.sel = q.sel.Where(sq.Gt{expiresAtColumnName: expiresAt.UTC()})
	q.del = q.del.Where(sq.Gt{expiresAtColumnName: expiresAt.UTC()})
	return q
}
