# Квази-сервис executor_fallback

Пример взаимодействия

```bash
python3 main.py
grpcurl -plaintext -d '{"display_name": "Natalia"}' localhost:9095 executor_profile.ExecutorProfileService/GetExecutorProfile
```
