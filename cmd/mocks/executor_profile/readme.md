# Заглушка сервиса executor_profile

Пример взаимодействия

```bash
python3 main.py
grpcurl -plaintext -d '{"display_name": "Natalia"}' localhost:9094 executor_profile.ExecutorProfileService/GetExecutorProfile
```
