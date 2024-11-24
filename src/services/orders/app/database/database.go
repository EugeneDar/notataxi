package database

import (
	"app/src/services/orders/app/model"
	"app/src/services/orders/app/utils"
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func EstablishConnection() (err error) {
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=%s",
		utils.GetenvSafe("POSTGRES_HOST"),
		utils.GetenvSafe("POSTGRES_PORT"),
		utils.GetenvSafe("POSTGRES_USER"),
		utils.GetenvSafe("POSTGRES_PASSWORD"),
		utils.GetenvSafe("POSTGRES_DB"),
		utils.GetenvSafe("TIME_ZONE"),
	)
	log.Printf("attempting to establish connection with database at %s\n", connectionString)
	db, err = sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}
	return nil
}

func AddAssignedOrder(assignedOrder *model.AssignedOrder) error {
	_, err := db.Exec(`
		UPDATE assigned_orders
		SET ExecutionStatus = 'completed'
		WHERE (ExecutionStatus = 'assigned' OR ExecutionStatus = 'acquired')
			AND ExecutorId = $1`,
		assignedOrder.ExecutorId)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
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
	return err
}

func AcquireAssignedOrder(executorId string) (*model.AssignedOrder, error) {
	row := db.QueryRow(`
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
	if err == sql.ErrNoRows {
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
	res, err := db.Exec(`
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
	affected, err := res.RowsAffected()
	return affected == 1, err
}

func CleanDatabase() error {
	_, err := db.Exec(`
        UPDATE assigned_orders 
        SET ExecutionStatus = 'completed' 
        WHERE ExecutionStatus IN ('assigned', 'acquired', 'cancelled')
    `)
	return err
}

func CleanTestOrders() error {
    _, err := db.Exec(`
        UPDATE assigned_orders 
        SET ExecutionStatus = 'cancelled' 
        WHERE (ExecutionStatus = 'assigned' OR ExecutionStatus = 'acquired')
            AND ExecutorId LIKE 'test_%'
    `)
    return err
}