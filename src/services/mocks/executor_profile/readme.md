```
python3 main.py
grpcurl -plaintext -d '{"display_name": "Natalia"}' localhost:50051 executor_profile.ExecutorProfileService/GetExecutorProfile
```
