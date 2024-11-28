# Примеры запросов

```bash
curl -X PUT "http://localhost:8080/order/assign?order_id=123&executor_id=456&zone_id=666"
curl -X GET "http://localhost:8080/order/acquire?executor_id=456"
curl -X POST "http://localhost:8080/order/cancel?order_id=123"
```
