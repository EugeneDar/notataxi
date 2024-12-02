# Tests on sources service

Before running the tests set the `MOCKS_ADDRESS` and `SOURCES_ADDRESS` environment variables. Example:

```bash
export MOCKS_ADDRESS=130.193.46.89
export MOCKS_ADDRESS=158.160.128.53
```

To launch functional tests run from the root directory of the repository:

```bash
python3 -m pytest test/sources/tests/test_sources_functional.py
```

To launch load tests run from the root directory of the repository:

```bash
locust -f test/sources/tests/test_sources_load.py
```
