```
python3 main.py
grpcurl -plaintext -d '{"zone_id": "zone_id_123"}' localhost:50051 zone_data.ZoneDataService/GetZoneData
```
