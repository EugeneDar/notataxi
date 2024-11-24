package database

import (
	"app/src/services/orders/app/utils"
	"app/src/services/sources/protobufs/sources"
	"database/sql"
	"fmt"
	"log"
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

func AddAssignedOrder(assignedOrder *sources.SourcesResponse) error {
	_, err := db.Exec(`
		UPDATE assigned_orders
		SET ExecutionStatus = 'completed'
		WHERE (ExecutionStatus = 'assigned' OR ExecutionStatus = 'acquired')
			AND ExecutorId = $1`,
		assignedOrder.GetExecutorProfile().GetId())
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
		"AssignedOrderId", // TODO generate UUID?
		assignedOrder.GetOrderId(),
		assignedOrder.GetExecutorProfile().GetId(),
		"ExecutionStatus", // TODO not enough info in sources service
		assignedOrder.GetPriceComponents().GetCoinCoeff(),
		assignedOrder.GetPriceComponents().GetBonusAmount(),
		assignedOrder.GetFinalCoinAmount(),
		assignedOrder.GetZoneDisplayName(),
		assignedOrder.GetUsedExecutorFallback(),
	)
	return err
}

func AcquireAssignedOrder(executorId string) (*sources.SourcesResponse, error) {
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
