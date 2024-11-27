import sys
import os
import pytest
import requests
import json
import uuid
from datetime import datetime
import time

ORDERS_ADDRESS = os.getenv("ORDERS_ADDRESS", "localhost:8080")
BASE_URL = f"http://{ORDERS_ADDRESS}"

def print_response_details(response):
    print(f"\nResponse Status Code: {response.status_code}")
    print(f"Response Headers: {dict(response.headers)}")
    try:
        print(f"Response Body: {json.dumps(response.json(), indent=2)}")
    except:
        print(f"Response Text: {response.text}")

def generate_unique_id():
    timestamp = datetime.now().strftime('%Y%m%d_%H%M%S')
    unique = uuid.uuid4().hex[:6]
    return f"test_{timestamp}_{unique}"

@pytest.fixture(scope="session", autouse=True)
def clean_test_orders():
    response = requests.post(f"{BASE_URL}/testing/clean-test-orders")
    print("\nCleaning test orders:")
    print_response_details(response)
    assert response.status_code == 200, f"Failed to clean test orders. Response: {response.text}"
    
    yield
    
    requests.post(f"{BASE_URL}/testing/clean-test-orders")

@pytest.fixture
def test_order():
    params = {
        "order_id": generate_unique_id(),
        "executor_id": "test_executor_456",
        "zone_id": "test_zone_789"
    }
    
    response = requests.put(f"{BASE_URL}/order/assign", params=params)
    assert response.status_code == 200
    
    yield params
    
    requests.post(f"{BASE_URL}/order/cancel", params={"order_id": params["order_id"]})

def test_assign_order():
    params = {
        "order_id": generate_unique_id(),
        "executor_id": "test_executor_456",
        "zone_id": "test_zone_789"
    }
    response = requests.put(f"{BASE_URL}/order/assign", params=params)
    print_response_details(response)
    assert response.status_code == 200
    assert response.json()["message"] == "Successfully created"

def test_assign_order_missing_params():
    response = requests.put(f"{BASE_URL}/order/assign")
    print_response_details(response)
    assert response.status_code == 400
    assert response.json()["message"] == "Missing parameters, please provide order_id, executor_id and zone_id"

def test_assign_order_duplicate(test_order):
    response = requests.put(f"{BASE_URL}/order/assign", params=test_order)
    print_response_details(response)
    assert response.status_code == 400
    assert response.json()["message"] == "AssignedOrder with provided orderId already exists"

def test_acquire_order(test_order):
    params = {"executor_id": test_order["executor_id"]}
    response = requests.get(f"{BASE_URL}/order/acquire", params=params)
    print_response_details(response)
    assert response.status_code == 200
    order = response.json()["order_profile"]
    assert order["assigned_order_id"] is not None
    assert order["order_id"] == test_order["order_id"]
    assert order["executor_id"] == test_order["executor_id"]
    assert order["execution_status"] == "acquired"

def test_acquire_order_no_orders():
    params = {"executor_id": "nonexistent_executor"}
    response = requests.get(f"{BASE_URL}/order/acquire", params=params)
    print_response_details(response)
    assert response.status_code == 200
    assert response.json()["message"] == "There are no orders assigned to you"

def test_cancel_order(test_order):
    params = {"order_id": test_order["order_id"]}
    response = requests.post(f"{BASE_URL}/order/cancel", params=params)
    print_response_details(response)
    assert response.status_code == 200
    assert response.json()["message"] == "Successfully canceled"

def test_cancel_nonexistent_order():
    params = {"order_id": "nonexistent_order"}
    response = requests.post(f"{BASE_URL}/order/cancel", params=params)
    print_response_details(response)
    assert response.status_code == 200
    assert "AssignedOrder with OrderId nonexistent_order does not exist" in response.json()["message"]

def test_cancel_missing_order_id():
    response = requests.post(f"{BASE_URL}/order/cancel")
    print_response_details(response)
    assert response.status_code == 400
    assert response.json()["message"] == "Missing parameters, please provide order_id"

def test_cancel_order_edge_cases():
    order_id = generate_unique_id()
    executor_id = f"test_executor_{generate_unique_id()}"
    
    assign_params = {
        "order_id": order_id,
        "executor_id": executor_id,
        "zone_id": "test_zone_789"
    }
    requests.put(f"{BASE_URL}/order/assign", params=assign_params)
    
    cancel_response = requests.post(f"{BASE_URL}/order/cancel", params={"order_id": order_id})
    assert cancel_response.status_code == 200
    
    cancel_response_2 = requests.post(f"{BASE_URL}/order/cancel", params={"order_id": order_id})
    assert cancel_response_2.status_code == 200
    assert "does not exist or has already been canceled" in cancel_response_2.json()["message"]
    
    acquire_response = requests.get(f"{BASE_URL}/order/acquire", params={"executor_id": executor_id})
    assert acquire_response.status_code == 200
    assert acquire_response.json()["message"] == "There are no orders assigned to you"
