# Tests on sources service

To launch functional tests run from the root directory of the repository:

```bash
python3 -m pytest test/sources/tests/test_sources_functional.py
```

To launch load tests run from the root directory of the repository:

```bash
locust -f test/sources/tests/test_sources_load.py
```
