package pg

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v4/stdlib" // driver
	"github.com/jmoiron/sqlx"
	"github.com/rendau/gms_temp/internal/interfaces"
	"time"
)

const ErrMsg = "PG-error"
const TransactionCtxKey = "pg_transaction"

type St struct {
	lg interfaces.Logger

	Db *sqlx.DB
}

type conSt interface {
	sqlx.ExtContext

	PrepareNamed(query string) (*sqlx.NamedStmt, error)
	GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

type txContainerSt struct {
	tx *sqlx.Tx
}

func NewSt(lg interfaces.Logger, dsn string) (*St, error) {
	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(100)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(10 * time.Minute)

	return &St{
		lg: lg,
		Db: db,
	}, nil
}

func (d *St) handleError(err error) error {
	if err == nil {
		return nil
	}

	// errStr := err.Error()

	d.lg.Errorw(ErrMsg, err)

	return err
}

func (d *St) getCon(ctx context.Context) conSt {
	if tx := d.getContextTransaction(ctx); tx != nil {
		return tx
	}
	return d.Db
}

func (d *St) getContextTransactionContainer(ctx context.Context) *txContainerSt {
	contextV := ctx.Value(TransactionCtxKey)
	if contextV == nil {
		return nil
	}

	switch tx := contextV.(type) {
	case *txContainerSt:
		return tx
	default:
		return nil
	}
}

func (d *St) getContextTransaction(ctx context.Context) *sqlx.Tx {
	container := d.getContextTransactionContainer(ctx)
	if container != nil {
		return container.tx
	}

	return nil
}

func (d *St) ContextWithTransaction(ctx context.Context) (context.Context, error) {
	tx, err := d.Db.BeginTxx(ctx, nil)
	if err != nil {
		return ctx, err
	}

	ctx = context.WithValue(ctx, TransactionCtxKey, &txContainerSt{tx: tx})

	return ctx, nil
}

func (d *St) CommitContextTransaction(ctx context.Context) error {
	tx := d.getContextTransaction(ctx)
	if tx == nil {
		return nil
	}

	if err := tx.Commit(); err != nil && err != sql.ErrTxDone && err != context.Canceled {
		d.lg.Errorw("Fail to commit transaction", err)
		return err
	}

	return nil
}

func (d *St) RollbackContextTransaction(ctx context.Context) {
	tx := d.getContextTransaction(ctx)
	if tx == nil {
		return
	}

	if err := tx.Rollback(); err != nil && err != sql.ErrTxDone && err != context.Canceled {
		d.lg.Errorw("Fail to rollback transaction", err)
	}
}

func (d *St) RenewContextTransaction(ctx context.Context) error {
	var err error

	container := d.getContextTransactionContainer(ctx)
	if container == nil {
		d.lg.Errorw(ErrMsg+": Transaction container not found in context", nil)
		return nil
	}

	if container.tx != nil {
		err = container.tx.Commit()
		if err != nil && err != sql.ErrTxDone {
			// try to rollback
			_ = container.tx.Rollback()

			d.lg.Errorw("Fail to commit transaction", err)
			return err
		}
	}

	container.tx, err = d.Db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}
