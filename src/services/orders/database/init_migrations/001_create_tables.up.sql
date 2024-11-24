SET TIMEZONE="Europe/Moscow";

CREATE TYPE order_execution_status AS ENUM ('assigned', 'acquired', 'cancelled', 'completed');

CREATE TABLE assigned_orders (
    AssignedOrderId VARCHAR(36) NOT NULL,  -- uuid as string
    OrderId VARCHAR(36) NOT NULL,
    ExecutorId VARCHAR(36) NOT NULL,
    ExecutionStatus order_execution_status NOT NULL,
    CoinCoefficient DOUBLE PRECISION NOT NULL,
    CoinBonusAmount DOUBLE PRECISION NOT NULL,
    FinalCoinAmount DOUBLE PRECISION NOT NULL,
    ZoneName VARCHAR(255),
    HasExecutorFallbackBeenUsed BOOLEAN NOT NULL,
    AssignTime TIMESTAMP WITH TIME ZONE NOT NULL,
    LastAcquireTime TIMESTAMP WITH TIME ZONE,

    PRIMARY KEY (AssignedOrderId),
    UNIQUE (OrderId)
);

CREATE UNIQUE INDEX assigned_orders_pk_total_search_index
    ON assigned_orders(AssignedOrderId);

CREATE UNIQUE INDEX assigned_orders_cancel_by_order_id_index
    ON assigned_orders(OrderId)
    WHERE ExecutionStatus = 'assigned' OR ExecutionStatus = 'acquired';

CREATE UNIQUE INDEX assigned_orders_search_by_executor_id_index
    ON assigned_orders(ExecutorId)
    WHERE ExecutionStatus = 'assigned' OR ExecutionStatus = 'acquired';
