```
python3 main.py
grpcurl -plaintext -d '{"order_id": "some_order_id"}' localhost:50051 order_data.OrderDataService/GetOrderData
```
