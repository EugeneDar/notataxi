package database

import (
	"notataxi/internal/orders/model"

	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
)

const (
	host        = "rc1b-xy6y7apt7j7jdpc7.mdb.yandexcloud.net,rc1d-opy4a78yulgzu7z2.mdb.yandexcloud.net"
	port        = 6432
	user        = "user1"
	password    = "NgdXRLUNn67d8tR"
	dbname      = "db1"
	sslrootcert = "/go/.postgresql/root.crt"
	timeZone    = "Europe/Moscow"
)

var db *pgxpool.Pool

func EstablishConnection() (err error) {
	connString := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=require sslrootcert=%s target_session_attrs=read-write",
		host, port, dbname, user, password, sslrootcert)

	db, err = pgxpool.New(context.Background(), connString)
	return err
}

func AddAssignedOrder(assignedOrder *model.AssignedOrder) (bool, error) {
	_, err := db.Exec(context.Background(), `
		UPDATE assigned_orders
		SET ExecutionStatus = 'completed'
		WHERE (ExecutionStatus = 'assigned' OR ExecutionStatus = 'acquired')
			AND ExecutorId = $1`,
		assignedOrder.ExecutorId)
	if err != nil {
		return false, err
	}
	_, err = db.Exec(context.Background(), `
		INSERT INTO assigned_orders (
				AssignedOrderId,
				OrderId,
				ExecutorId,
				ExecutionStatus,
				CoinCoefficient,
				CoinBonusAmount,
				FinalCoinAmount,
				ZoneName,
				HasExecutorFallbackBeenUsed,
				AssignTime,
				LastAcquireTime
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NULL)`,
		assignedOrder.AssignedOrderId,
		assignedOrder.OrderId,
		assignedOrder.ExecutorId,
		assignedOrder.ExecutionStatus,
		assignedOrder.CoinCoefficient,
		assignedOrder.CoinBonusAmount,
		assignedOrder.FinalCoinAmount,
		assignedOrder.ZoneName,
		assignedOrder.HasExecutorFallbackBeenUsed,
	)
	if err != nil && err.Error() == "ERROR: duplicate key value violates unique constraint \"assigned_orders_orderid_key\" (SQLSTATE 23505)" {
		return false, nil
	}
	return true, err
}

func AcquireAssignedOrder(executorId string) (*model.AssignedOrder, error) {
	row := db.QueryRow(context.Background(), `
		UPDATE assigned_orders
		SET
			ExecutionStatus = 'acquired',
			LastAcquireTime = NOW()
		WHERE
			(ExecutionStatus = 'assigned' OR ExecutionStatus = 'acquired')
			AND ExecutorId = $1
		RETURNING
			AssignedOrderId,
			OrderId,
			ExecutorId,
			ExecutionStatus,
			CoinCoefficient,
			CoinBonusAmount,
			FinalCoinAmount,
			ZoneName,
			HasExecutorFallbackBeenUsed,
			AssignTime,
			LastAcquireTime`,
		executorId,
	)
	var assignedOrder model.AssignedOrder
	var zoneName sql.NullString
	var lastAcquireTime sql.NullTime
	err := row.Scan(
		&assignedOrder.AssignedOrderId,
		&assignedOrder.OrderId,
		&assignedOrder.ExecutorId,
		&assignedOrder.ExecutionStatus,
		&assignedOrder.CoinCoefficient,
		&assignedOrder.CoinBonusAmount,
		&assignedOrder.FinalCoinAmount,
		&zoneName,
		&assignedOrder.HasExecutorFallbackBeenUsed,
		&assignedOrder.AssignTime,
		&lastAcquireTime,
	)
	if err != nil && err.Error() == "no rows in result set" {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if zoneName.Valid {
		assignedOrder.ZoneName = zoneName.String
	} else {
		assignedOrder.ZoneName = "unknown"
	}
	if lastAcquireTime.Valid {
		assignedOrder.LastAcquireTime = lastAcquireTime.Time
	} else {
		return nil, errors.New("unexpected unset AcquireTime during acquiring order")
	}
	return &assignedOrder, nil
}

func CancelAssignedOrder(orderId string) (bool, error) {
	res, err := db.Exec(context.Background(), `
		UPDATE assigned_orders
		SET ExecutionStatus = 'cancelled'
		WHERE (ExecutionStatus = 'assigned' OR ExecutionStatus = 'acquired')
			AND OrderId = $1
		RETURNING 1`,
		orderId,
	)
	if err != nil {
		return false, err
	}
	affected := res.RowsAffected()
	return affected == 1, nil
}

func CleanDatabase() error {
	_, err := db.Exec(context.Background(), `
        UPDATE assigned_orders 
        SET ExecutionStatus = 'completed' 
        WHERE ExecutionStatus IN ('assigned', 'acquired', 'cancelled')
    `)
	return err
}

func CleanTestOrders() error {
	_, err := db.Exec(context.Background(), `
        UPDATE assigned_orders 
        SET ExecutionStatus = 'cancelled' 
        WHERE (ExecutionStatus = 'assigned' OR ExecutionStatus = 'acquired')
            AND ExecutorId LIKE 'test_%'
    `)
	return err
}
