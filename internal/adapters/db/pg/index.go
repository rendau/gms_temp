package pg

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	_ "github.com/jackc/pgx/v4/stdlib" // driver
	"github.com/rendau/gms_temp/internal/domain/errs"
	"github.com/rendau/gms_temp/internal/interfaces"
	"strconv"
	"strings"
	"time"
)

const ErrMsg = "PG-error"
const TransactionCtxKey = "pg_transaction"

type St struct {
	lg interfaces.Logger

	Con *pgxpool.Pool
}

type conSt interface {
	Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

type txContainerSt struct {
	tx pgx.Tx
}

func NewSt(lg interfaces.Logger, dsn string) (*St, error) {
	dbConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		lg.Errorw(ErrMsg+": Fail to parse dsn", err)
		return nil, err
	}

	dbConfig.MaxConns = 30
	dbConfig.MinConns = 3
	dbConfig.MaxConnLifetime = 0
	dbConfig.MaxConnIdleTime = 3 * time.Minute
	dbConfig.HealthCheckPeriod = 20 * time.Second
	dbConfig.LazyConnect = true

	dbPool, err := pgxpool.ConnectConfig(context.Background(), dbConfig)
	if err != nil {
		lg.Errorw(ErrMsg+": Fail to connect to db", err)
		return nil, err
	}

	return &St{
		lg:  lg,
		Con: dbPool,
	}, nil
}

func (d *St) handleError(err error) error {
	if err == nil {
		return nil
	}

	if err == context.Canceled ||
		err == context.DeadlineExceeded {
		d.lg.Infow("PG-error, context canceled")
		return errs.ServiceNA
	}

	d.lg.Errorw(ErrMsg, err)

	return err
}

func (d *St) getCon(ctx context.Context) conSt {
	if tx := d.getContextTransaction(ctx); tx != nil {
		return tx
	}
	return d.Con
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

func (d *St) getContextTransaction(ctx context.Context) pgx.Tx {
	container := d.getContextTransactionContainer(ctx)
	if container != nil {
		return container.tx
	}

	return nil
}

func (d *St) ContextWithTransaction(ctx context.Context) (context.Context, error) {
	tx, err := d.Con.Begin(ctx)
	if err != nil {
		return ctx, d.handleError(err)
	}

	return context.WithValue(ctx, TransactionCtxKey, &txContainerSt{tx: tx}), nil
}

func (d *St) CommitContextTransaction(ctx context.Context) error {
	tx := d.getContextTransaction(ctx)
	if tx == nil {
		return nil
	}

	err := tx.Commit(ctx)
	if err != nil {
		if err == context.Canceled ||
			err == context.DeadlineExceeded {
			return errs.ContextCancelled
		}
		if err != pgx.ErrTxClosed &&
			err != pgx.ErrTxCommitRollback {
			return d.handleError(err)
		}
	}

	return nil
}

func (d *St) RollbackContextTransaction(ctx context.Context) {
	tx := d.getContextTransaction(ctx)
	if tx == nil {
		return
	}

	err := tx.Rollback(ctx)
	if err != nil {
		if err == context.Canceled ||
			err == context.DeadlineExceeded {
			return
		}
		if err != pgx.ErrTxClosed &&
			err != pgx.ErrTxCommitRollback {
			_ = d.handleError(err)
		}
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
		err = container.tx.Commit(ctx)
		if err != nil &&
			err != pgx.ErrTxClosed &&
			err != pgx.ErrTxCommitRollback &&
			err != context.Canceled &&
			err != context.DeadlineExceeded {
			// try to rollback
			_ = container.tx.Rollback(ctx)
			return d.handleError(err)
		}
	}

	container.tx, err = d.Con.Begin(ctx)
	if err != nil {
		return d.handleError(err)
	}

	return nil
}

func (d *St) dbExec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	return d.getCon(ctx).Exec(ctx, sql, args...)
}

func (d *St) dbQuery(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return d.getCon(ctx).Query(ctx, sql, args...)
}

func (d *St) dbQueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	return d.getCon(ctx).QueryRow(ctx, sql, args...)
}

func (d *St) queryRebindNamed(sql string, argMap map[string]interface{}) (string, []interface{}) {
	resultQuery := sql
	args := make([]interface{}, 0, len(argMap))

	for k, v := range argMap {
		if strings.Contains(resultQuery, "${"+k+"}") {
			args = append(args, v)
			resultQuery = strings.ReplaceAll(resultQuery, "${"+k+"}", "$"+strconv.Itoa(len(args)))
		}
	}

	return resultQuery, args
}

func (d *St) dbExecM(ctx context.Context, sql string, argMap map[string]interface{}) (pgconn.CommandTag, error) {
	rbSql, args := d.queryRebindNamed(sql, argMap)

	return d.getCon(ctx).Exec(ctx, rbSql, args...)
}

func (d *St) dbQueryM(ctx context.Context, sql string, argMap map[string]interface{}) (pgx.Rows, error) {
	rbSql, args := d.queryRebindNamed(sql, argMap)

	return d.getCon(ctx).Query(ctx, rbSql, args...)
}

func (d *St) dbQueryRowM(ctx context.Context, sql string, argMap map[string]interface{}) pgx.Row {
	rbSql, args := d.queryRebindNamed(sql, argMap)

	return d.getCon(ctx).QueryRow(ctx, rbSql, args...)
}
