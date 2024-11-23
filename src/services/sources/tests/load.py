import sys
import os
import grpc
import time
import pytest
import concurrent.futures
from google.protobuf import empty_pb2
sys.path.append(os.path.join(os.path.dirname(os.path.abspath(__file__)), '../protobufs'))
from config_pb2_grpc import ConfigServiceStub
from executor_profile_pb2_grpc import ExecutorProfileServiceStub
from executor_profile_pb2 import ExecutorProfileRequest
from order_data_pb2_grpc import OrderDataServiceStub
from order_data_pb2 import OrderDataRequest
from sources_pb2_grpc import OrderInfoServiceStub
from sources_pb2 import OrderInfoRequest
from toll_roads_pb2_grpc import TollRoadsServiceStub
from toll_roads_pb2 import TollRoadsRequest
from zone_data_pb2_grpc import ZoneDataServiceStub
from zone_data_pb2 import ZoneDataRequest

@pytest.fixture(scope="module")
def grpc_channel():
    with grpc.insecure_channel('localhost:50051') as channel:
        yield channel

@pytest.fixture(scope="module")
def config_service(grpc_channel):
    return ConfigServiceStub(grpc_channel)

@pytest.fixture(scope="module")
def executor_profile_service(grpc_channel):
    return ExecutorProfileServiceStub(grpc_channel)

@pytest.fixture(scope="module")
def order_data_service(grpc_channel):
    return OrderDataServiceStub(grpc_channel)

@pytest.fixture(scope="module")
def order_info_service(grpc_channel):
    return OrderInfoServiceStub(grpc_channel)

@pytest.fixture(scope="module")
def toll_roads_service(grpc_channel):
    return TollRoadsServiceStub(grpc_channel)

@pytest.fixture(scope="module")
def zone_data_service(grpc_channel):
    return ZoneDataServiceStub(grpc_channel)

CONCURRENT_REQUESTS = 100

def load_test_function(service_call, request, expected_code=grpc.StatusCode.OK):
    try:
        start = time.time()
        response = service_call(request)
        elapsed = time.time() - start
        assert isinstance(response, expected_code)
        return elapsed
    except grpc.RpcError as e:
        assert e.code() == expected_code, f"Expected {expected_code} but got {e.code()}"
        return None

def test_load_config_service(config_service):
    with concurrent.futures.ThreadPoolExecutor(max_workers=CONCURRENT_REQUESTS) as executor:
        futures = [
            executor.submit(load_test_function, config_service.GetConfig, empty_pb2.Empty())
            for _ in range(CONCURRENT_REQUESTS)
        ]
        results = [future.result() for future in futures if future.result() is not None]
        print("ConfigService - Average response time:", sum(results) / len(results))

def test_load_executor_profile_service(executor_profile_service):
    request = ExecutorProfileRequest(display_name="Test User")
    with concurrent.futures.ThreadPoolExecutor(max_workers=CONCURRENT_REQUESTS) as executor:
        futures = [
            executor.submit(load_test_function, executor_profile_service.GetExecutorProfile, request)
            for _ in range(CONCURRENT_REQUESTS)
        ]
        results = [future.result() for future in futures if future.result() is not None]
        print("ExecutorProfileService - Average response time:", sum(results) / len(results))

def test_load_order_data_service(order_data_service):
    request = OrderDataRequest(order_id="12345")
    with concurrent.futures.ThreadPoolExecutor(max_workers=CONCURRENT_REQUESTS) as executor:
        futures = [
            executor.submit(load_test_function, order_data_service.GetOrderData, request)
            for _ in range(CONCURRENT_REQUESTS)
        ]
        results = [future.result() for future in futures if future.result() is not None]
        print("OrderDataService - Average response time:", sum(results) / len(results))

def test_load_order_info_service(order_info_service):
    request = OrderInfoRequest(order_id="12345", executor_id="executor_01")
    with concurrent.futures.ThreadPoolExecutor(max_workers=CONCURRENT_REQUESTS) as executor:
        futures = [
            executor.submit(load_test_function, order_info_service.GetOrderInfo, request)
            for _ in range(CONCURRENT_REQUESTS)
        ]
        results = [future.result() for future in futures if future.result() is not None]
        print("OrderInfoService - Average response time:", sum(results) / len(results))

def test_load_toll_roads_service(toll_roads_service):
    request = TollRoadsRequest(display_name="Highway_101")
    with concurrent.futures.ThreadPoolExecutor(max_workers=CONCURRENT_REQUESTS) as executor:
        futures = [
            executor.submit(load_test_function, toll_roads_service.GetTollRoads, request)
            for _ in range(CONCURRENT_REQUESTS)
        ]
        results = [future.result() for future in futures if future.result() is not None]
        print("TollRoadsService - Average response time:", sum(results) / len(results))

def test_load_zone_data_service(zone_data_service):
    request = ZoneDataRequest(zone_id="zone_123")
    with concurrent.futures.ThreadPoolExecutor(max_workers=CONCURRENT_REQUESTS) as executor:
        futures = [
            executor.submit(load_test_function, zone_data_service.GetZoneData, request)
            for _ in range(CONCURRENT_REQUESTS)
        ]
        results = [future.result() for future in futures if future.result() is not None]
        print("ZoneDataService - Average response time:", sum(results) / len(results))
