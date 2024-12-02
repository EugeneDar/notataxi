# Tests on orders service

Before running the tests set the `ORDERS_ADDRESS` environment variable. Example:

```bash
export ORDERS_ADDRESS=158.160.140.176
```

To launch functional tests run from the root directory of the repository:

```bash
python3 -m pytest test/orders/tests/test_orders_functional.py
```

To launch load tests run from the root directory of the repository:

```bash
locust -f test/orders/tests/test_orders_load.py --host=http://$ORDERS_ADDRESS:8080
```
