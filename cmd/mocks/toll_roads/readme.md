```
python3 main.py
grpcurl -plaintext -d '{"display_name": "zone_display_name"}' localhost:9093 toll_roads.TollRoadsService/GetTollRoads
```
