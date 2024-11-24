```
python3 main.py
grpcurl -plaintext -d '{"display_name": "Natalia"}' localhost:9094 executor_profile.ExecutorProfileService/GetExecutorProfile
```
