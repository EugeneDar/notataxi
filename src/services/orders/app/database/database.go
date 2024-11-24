package database

import (
	"app/src/services/sources/protobufs/sources"
	"database/sql"
	"fmt"
	"os"

	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v4"
)

const (
	host     = "rc1b-xy6y7apt7j7jdpc7.mdb.yandexcloud.net,rc1d-opy4a78yulgzu7z2.mdb.yandexcloud.net"
	port     = 6432
	user     = "user1"
	password = "NgdXRLUNn67d8tR"
	dbname   = "db1"
	ca       = "/root/.postgresql/root.crt"
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
		fmt.Fprintf(os.Stderr, "Unable to parse config: %v\n", err)
		os.Exit(1)
	}

	connConfig.TLSConfig = &tls.Config{
		RootCAs:            rootCertPool,
		InsecureSkipVerify: true,
	}

	db, err = pgx.ConnectConfig(context.Background(), connConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return nil
}

func AddAssignedOrder(assignedOrder *sources.SourcesResponse) error {
	_, err := db.Exec(context.Background(), `
		UPDATE assigned_orders
		SET ExecutionStatus = 'completed'
		WHERE (ExecutionStatus = 'assigned' OR ExecutionStatus = 'acquired')
			AND ExecutorId = $1`,
		assignedOrder.GetExecutorProfile().GetId())
	if err != nil {
		return err
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
		uuid.New().String(),
		assignedOrder.GetOrderId(),
		assignedOrder.GetExecutorProfile().GetId(),
		"assigned",
		assignedOrder.GetPriceComponents().GetCoinCoeff(),
		assignedOrder.GetPriceComponents().GetBonusAmount(),
		assignedOrder.GetFinalCoinAmount(),
		assignedOrder.GetZoneDisplayName(),
		assignedOrder.GetUsedExecutorFallback(),
	)
	return err
}

func AcquireAssignedOrder(executorId string) (*sources.SourcesResponse, error) {
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
	var assignedOrder sources.SourcesResponse
	var zoneName sql.NullString

	assignedOrderId := ""
	status := ""

	err := row.Scan(
		&assignedOrderId,
		&assignedOrder.OrderId,
		&assignedOrder.ExecutorProfile.Id,
		&status,
		&assignedOrder.GetPriceComponents().CoinCoeff,
		&assignedOrder.GetPriceComponents().BonusAmount,
		&assignedOrder.FinalCoinAmount,
		&zoneName,
		&assignedOrder.UsedExecutorFallback,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if zoneName.Valid {
		assignedOrder.ZoneDisplayName = zoneName.String
	} else {
		assignedOrder.ZoneDisplayName = "unknown"
	}

	return &assignedOrder, nil
}

func CancelAssignedOrder(orderId string) bool {
	res, err := db.Exec(context.Background(), `
		UPDATE assigned_orders
		SET ExecutionStatus = 'cancelled'
		WHERE (ExecutionStatus = 'assigned' OR ExecutionStatus = 'acquired')
			AND OrderId = $1
		RETURNING 1`,
		orderId,
	)
	if err != nil {
		return false
	}
	affected := res.RowsAffected()
	return affected == 1
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
