# Заглушка сервиса order_data

Пример взаимодействия

```bash
python3 main.py
grpcurl -plaintext -d '{"order_id": "some_order_id"}' localhost:9091 order_data.OrderDataService/GetOrderData
```
