```
python3 main.py
grpcurl -plaintext -d '{"display_name": "zone_display_name"}' localhost:50051 toll_roads.TollRoadsService/GetTollRoads
```
