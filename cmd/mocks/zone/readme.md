# Заглушка сервиса zone

Пример взаимодействия

```bash
python3 main.py
grpcurl -plaintext -d '{"zone_id": "zone_id_123"}' localhost:9092 zone_data.ZoneDataService/GetZoneData
```
