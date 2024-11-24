import os
import sys
from locust import User, task, between, events, runners
import grpc
from google.protobuf import empty_pb2
import time
import gevent
import random
from datetime import datetime
import statistics

sys.path.append(os.path.join(os.path.dirname(os.path.abspath(__file__)), '../protobufs'))

from sources_pb2_grpc import SourcesServiceStub
from sources_pb2 import SourcesRequest
from executor_profile_pb2 import ExecutorProfileRequest
from order_data_pb2 import OrderDataRequest
from zone_data_pb2 import ZoneDataRequest
from toll_roads_pb2 import TollRoadsRequest

GRPC_HOST = "localhost:9000"

class GRPCUser(User):
    abstract = True
    wait_time = between(0.1, 1) 
    
    def __init__(self, environment):
        super().__init__(environment)
        self.client = None
        self.connection_successful = False
        
    def on_start(self):
        try:
            for attempt in range(3):
                try:
                    print(f"Attempting to connect to {GRPC_HOST} (attempt {attempt + 1}/3)")
                    self.channel = grpc.insecure_channel(GRPC_HOST)
                    # Уменьшаем timeout для быстрого фидбека
                    grpc.channel_ready_future(self.channel).result(timeout=2)
                    self.client = SourcesServiceStub(self.channel)
                    self.connection_successful = True
                    print(f"Successfully connected to gRPC server at {GRPC_HOST}")
                    break
                except grpc.FutureTimeoutError:
                    print(f"Connection attempt {attempt + 1} failed to {GRPC_HOST}")
                    if attempt < 2:  # Don't sleep on last attempt
                        time.sleep(1)
            
            if not self.connection_successful:
                print(f"Failed to connect to gRPC server at {GRPC_HOST} after 3 attempts")
                if isinstance(self.environment.runner, runners.LocalRunner):
                    self.environment.runner.quit()
                return

            self.order_ids = [f"order_{i}" for i in range(1000)]
            self.executor_ids = [f"exec_{i}" for i in range(100)]
            
        except Exception as e:
            print(f"Failed to initialize gRPC client: {str(e)}")
            if isinstance(self.environment.runner, runners.LocalRunner):
                self.environment.runner.quit()
        
    def on_stop(self):
        if hasattr(self, 'channel') and self.channel:
            try:
                self.channel.close()
            except Exception as e:
                print(f"Error closing channel: {e}")

class SourcesUser(GRPCUser):
    def make_request(self, name, request_func):
        if not self.connection_successful:
            print("Skipping request - no connection to server")
            return

        start_time = time.time()
        try:
            request_func()
            total_time = int((time.time() - start_time) * 1000)
            events.request.fire(
                request_type="grpc",
                name=name,
                response_time=total_time,
                response_length=0,
                exception=None,
            )
        except grpc.RpcError as e:
            total_time = int((time.time() - start_time) * 1000)
            events.request.fire(
                request_type="grpc",
                name=name,
                response_time=total_time,
                response_length=0,
                exception=e,
            )
            print(f"gRPC error during {name}: {e.code()} - {e.details()}")
            if e.code() == grpc.StatusCode.UNAVAILABLE:
                print(f"Service unavailable: {e.details()}")
                if isinstance(self.environment.runner, runners.LocalRunner):
                    self.environment.runner.quit()

    @task(1)
    def get_order_info(self):
        order_id = random.choice(self.order_ids)
        executor_id = random.choice(self.executor_ids)
        request = SourcesRequest(
            order_id=order_id,
            executor_id=executor_id
        )
        
        def request_func():
            try:
                response = self.client.GetOrderInfo(request, timeout=5)
                if not response.order_id:
                    raise ValueError("Empty response")
            except grpc.RpcError as e:
                if e.code() == grpc.StatusCode.DEADLINE_EXCEEDED:
                    print(f"Request timed out for order_id={order_id}")
                raise
                
        self.make_request("GetOrderInfo", request_func)

class LoadTestStats:
    def __init__(self):
        self.response_times = []
        self.requests = 0
        self.failures = 0
        self.start_time = time.time()
        
    @events.request.add_listener
    def on_request(self, request_type, name, response_time, response_length, exception, **kwargs):
        self.requests += 1
        if exception:
            self.failures += 1
            print(f"Request failed: {type(exception).__name__} - {str(exception)}")
        self.response_times.append(response_time)
        
    def get_stats(self):
        duration = time.time() - self.start_time
        if not self.response_times:
            return "No requests made"
            
        return {
            "total_requests": self.requests,
            "total_failures": self.failures,
            "avg_response_time_ms": statistics.mean(self.response_times),
            "median_response_time_ms": statistics.median(self.response_times),
            "95th_percentile_ms": statistics.quantiles(self.response_times, n=20)[-1],
            "requests_per_second": self.requests / duration,
            "failure_percentage": (self.failures / self.requests * 100) if self.requests > 0 else 0
        }

stats = LoadTestStats()

@events.test_start.add_listener
def on_test_start(environment, **kwargs):
    print(f"Test is starting... Connecting to {GRPC_HOST}")

@events.test_stop.add_listener
def on_test_stop(environment, **kwargs):
    print("\nTest is ending...")
    print_stats()

@events.quitting.add_listener
def print_final_stats(environment, **kwargs):
    print_stats()

def print_stats():
    final_stats = stats.get_stats()
    if isinstance(final_stats, dict):
        print("\n=== Load Test Results ===")
        print(f"Total Requests: {final_stats['total_requests']}")
        print(f"Total Failures: {final_stats['total_failures']}")
        print(f"Average Response Time: {final_stats['avg_response_time_ms']:.2f}ms")
        print(f"Median Response Time: {final_stats['median_response_time_ms']:.2f}ms")
        print(f"95th Percentile Response Time: {final_stats['95th_percentile_ms']:.2f}ms")
        print(f"Requests Per Second: {final_stats['requests_per_second']:.2f}")
        print(f"Failure Percentage: {final_stats['failure_percentage']:.2f}%")
    else:
        print(final_stats)

if __name__ == "__main__":
    pass