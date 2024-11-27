package database

import (
	"app/src/services/orders/model"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"

	"context"

	"github.com/jackc/pgx/v4"
)

const (
	host     = "rc1b-xy6y7apt7j7jdpc7.mdb.yandexcloud.net,rc1d-opy4a78yulgzu7z2.mdb.yandexcloud.net"
	port     = 6432
	user     = "user1"
	password = "NgdXRLUNn67d8tR"
	dbname   = "db1"
	ca       = "/go/src/services/orders/database/root.crt"
	timeZone = "Europe/Moscow"
)

var db *pgx.Conn

func EstablishConnection() (err error) {
	rootCertPool := x509.NewCertPool()
	pem, err := ioutil.ReadFile(ca)
	if err != nil {
		panic(err)
	}

	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		panic("Failed to append PEM.")
	}

	connstring := fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=verify-full target_session_attrs=read-write",
		host, port, dbname, user, password)

	connConfig, err := pgx.ParseConfig(connstring)
	if err != nil {
		return fmt.Errorf("unable to parse config: %v", err)
	}

	connConfig.TLSConfig = &tls.Config{
		RootCAs:            rootCertPool,
		InsecureSkipVerify: true,
	}

	db, err = pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}
	return nil
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
