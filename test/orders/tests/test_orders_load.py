from locust import HttpUser, task, between
from datetime import datetime
import uuid
import random

def generate_unique_id():
    """Generate a unique test ID combining timestamp and UUID"""
    timestamp = datetime.now().strftime('%Y%m%d_%H%M%S')
    unique = uuid.uuid4().hex[:6]
    return f"test_{timestamp}_{unique}"

class OrdersUser(HttpUser):
    wait_time = between(1, 3)
    
    def on_start(self):
        """Initialize user-specific data when the simulation starts"""
        self.executor_id = f"test_executor_{generate_unique_id()}"
        self.zone_id = "test_zone_789"
        self.active_orders = set()

    @task(3)
    def assign_order(self):
        """Test order assignment with unique order IDs"""
        order_id = generate_unique_id()
        params = {
            "order_id": order_id,
            "executor_id": self.executor_id,
            "zone_id": self.zone_id
        }
        
        with self.client.put("/order/assign", params=params, catch_response=True) as response:
            if response.status_code == 200:
                self.active_orders.add(order_id)
                response.success()
            elif response.status_code == 400 and "already exists" in response.text:
                response.success()
            else:
                response.failure(f"Assignment failed with status {response.status_code}")

    @task(2)
    def acquire_order(self):
        """Test order acquisition for the current executor"""
        params = {"executor_id": self.executor_id}
        
        with self.client.get("/order/acquire", params=params, catch_response=True) as response:
            if response.status_code == 200:
                if "There are no orders" not in response.text:
                    try:
                        order_info = response.json().get("order_profile", {})
                        if order_info.get("order_id"):
                            self.active_orders.add(order_info["order_id"])
                    except Exception:
                        pass
                response.success()
            else:
                response.failure(f"Acquisition failed with status {response.status_code}")

    @task(1)
    def cancel_order(self):
        """Test order cancellation for previously created orders"""
        if not self.active_orders:
            return
            
        order_id = random.choice(list(self.active_orders))
        params = {"order_id": order_id}
        
        with self.client.post("/order/cancel", params=params, catch_response=True) as response:
            if response.status_code == 200:
                self.active_orders.discard(order_id)
                response.success()
            elif "does not exist" in response.text:
                self.active_orders.discard(order_id)
                response.success()
            else:
                response.failure(f"Cancellation failed with status {response.status_code}")

class OrdersChaosUser(HttpUser):
    """Simulate chaos testing scenarios"""
    wait_time = between(0.1, 1)
    
    @task(3)
    def assign_invalid_orders(self):
        """Test with missing or invalid parameters"""
        scenarios = [
            {},
            {"order_id": generate_unique_id()},
            {"executor_id": "test_executor", "zone_id": "test_zone"},
            {"order_id": "", "executor_id": "", "zone_id": ""}
        ]
        
        scenario = random.choice(scenarios)
        self.client.put("/order/assign", params=scenario)

    @task(2)
    def rapid_acquire_attempts(self):
        """Rapidly attempt to acquire orders for random executors"""
        params = {"executor_id": f"chaos_executor_{uuid.uuid4().hex[:6]}"}
        self.client.get("/order/acquire", params=params)

    @task(1)
    def cancel_nonexistent_orders(self):
        """Attempt to cancel random or nonexistent orders"""
        params = {"order_id": f"nonexistent_{uuid.uuid4().hex}"}
        self.client.post("/order/cancel", params=params)
