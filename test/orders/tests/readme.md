# Tests on orders service

To launch functional tests run from the root directory of the repository:

```bash
python3 -m pytest test/orders/tests/test_orders_functional.py
```

To launch load tests run from the root directory of the repository:

```bash
locust -f test/orders/tests/test_orders_load.py --host=http://$ORDERS_ADDRESS:8080
```
