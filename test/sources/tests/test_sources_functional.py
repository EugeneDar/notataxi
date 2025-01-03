import sys
import os
import pytest
import grpc
from google.protobuf import empty_pb2

sys.path.append('internal/protobufs')
from internal.protobufs.config_pb2_grpc import ConfigServiceStub
from internal.protobufs.config_pb2 import ConfigResponse
from internal.protobufs.executor_profile_pb2_grpc import ExecutorProfileServiceStub
from internal.protobufs.executor_profile_pb2 import ExecutorProfileRequest
from internal.protobufs.order_data_pb2_grpc import OrderDataServiceStub
from internal.protobufs.order_data_pb2 import OrderDataRequest
from internal.protobufs.sources_pb2_grpc import SourcesServiceStub
from internal.protobufs.sources_pb2 import SourcesRequest
from internal.protobufs.toll_roads_pb2_grpc import TollRoadsServiceStub
from internal.protobufs.toll_roads_pb2 import TollRoadsRequest
from internal.protobufs.zone_data_pb2_grpc import ZoneDataServiceStub
from internal.protobufs.zone_data_pb2 import ZoneDataRequest

MOCKS_ADDRESS = os.getenv("MOCKS_ADDRESS")
if MOCKS_ADDRESS is None:
    sys.exit('Fill MOCKS_ADDRESS before running the test')

SOURCES_ADDRESS = os.getenv("SOURCES_ADDRESS")
if SOURCES_ADDRESS is None:
    sys.exit('Fill SOURCES_ADDRESS before running the test')

@pytest.fixture(scope="module")
def grpc_channel_config():
    with grpc.insecure_channel(f'{MOCKS_ADDRESS}:9090') as channel:
        yield channel

@pytest.fixture(scope="module")
def grpc_channel_executor_profile():
    with grpc.insecure_channel(f'{MOCKS_ADDRESS}:9094') as channel:
        yield channel

@pytest.fixture(scope="module")
def grpc_channel_order_data():
    with grpc.insecure_channel(f'{MOCKS_ADDRESS}:9091') as channel:
        yield channel

@pytest.fixture(scope="module")
def grpc_channel_sources():
    with grpc.insecure_channel(f'{SOURCES_ADDRESS}:9000') as channel:
        yield channel

@pytest.fixture(scope="module")
def grpc_channel_toll_roads():
    with grpc.insecure_channel(f'{MOCKS_ADDRESS}:9093') as channel:
        yield channel

@pytest.fixture(scope="module")
def grpc_channel_zone_data():
    with grpc.insecure_channel(f'{MOCKS_ADDRESS}:9092') as channel:
        yield channel

@pytest.fixture(scope="module")
def config_service(grpc_channel_config):
    return ConfigServiceStub(grpc_channel_config)

@pytest.fixture(scope="module")
def executor_profile_service(grpc_channel_executor_profile):
    return ExecutorProfileServiceStub(grpc_channel_executor_profile)

@pytest.fixture(scope="module")
def order_data_service(grpc_channel_order_data):
    return OrderDataServiceStub(grpc_channel_order_data)

@pytest.fixture(scope="module")
def sources_service(grpc_channel_sources):
    return SourcesServiceStub(grpc_channel_sources)

@pytest.fixture(scope="module")
def toll_roads_service(grpc_channel_toll_roads):
    return TollRoadsServiceStub(grpc_channel_toll_roads)

@pytest.fixture(scope="module")
def zone_data_service(grpc_channel_zone_data):
    return ZoneDataServiceStub(grpc_channel_zone_data)

def test_get_config(config_service):
    response = config_service.GetConfig(empty_pb2.Empty())
    assert isinstance(response.min_price, int)
    assert 1 <= response.min_price <= 1000, "Min price out of range"

def test_get_executor_profile_valid_request(executor_profile_service):
    request = ExecutorProfileRequest(executor_id="473831424")
    response = executor_profile_service.GetExecutorProfile(request)
    assert response.id
    assert len(response.id) > 0, "ID should not be empty"
    tags = list(response.tags)
    assert len(tags) > 0 and all([lambda x: isinstance(x, str) for x in tags])
    assert all(isinstance(tag, str) for tag in response.tags), "All tags should be strings"
    assert 0 <= response.rating <= 5, "Rating out of range"

def test_get_executor_profile_empty_name(executor_profile_service):
    request = ExecutorProfileRequest(executor_id="")
    with pytest.raises(grpc.RpcError) as exc_info:
        executor_profile_service.GetExecutorProfile(request)
    assert exc_info.value.code() == grpc.StatusCode.INVALID_ARGUMENT, "Should return INVALID_ARGUMENT for empty order"

def test_get_order_data_valid_order(order_data_service):
    request = OrderDataRequest(order_id="order_123")
    response = order_data_service.GetOrderData(request)
    assert response.order_id == "order_123"
    assert len(response.user_id) > 0, "User ID should not be empty"
    assert len(response.zone_id) > 0, "Zone ID should not be empty"
    assert response.base_coin_amount > 0, "Base coin amount should be positive"

def test_get_order_data_error_order(order_data_service):
    request = OrderDataRequest(order_id="")
    with pytest.raises(grpc.RpcError) as exc_info:
        order_data_service.GetOrderData(request)
    assert exc_info.value.code() == grpc.StatusCode.INVALID_ARGUMENT, "Should return INVALID_ARGUMENT for empty order"

def test_get_sources_valid_order(sources_service):
    request = SourcesRequest(order_id="order_123", executor_id="exec_456")
    response = sources_service.GetOrderInfo(request)
    assert response.order_id == "order_123"
    assert response.final_coin_amount > 0, "Final coin amount should be positive"
    assert response.price_components.base_coin_amount > 0, "Base coin amount should be positive"
    assert response.price_components.coin_coeff >= 1, "Coin coefficient should be >= 1"
    assert isinstance(response.executor_profile.rating, float)
    assert 0 <= response.executor_profile.rating <= 5, "Rating out of range"

def test_get_toll_roads_valid_request(toll_roads_service):
    request = TollRoadsRequest(display_name="Highway 101")
    response = toll_roads_service.GetTollRoads(request)
    assert response.bonus_amount >= 0, "Bonus amount should not be negative"

def test_get_toll_roads_invalid_request(toll_roads_service):
    request = TollRoadsRequest(display_name="")
    with pytest.raises(grpc.RpcError) as exc_info:
        toll_roads_service.GetTollRoads(request)
    assert exc_info.value.code() == grpc.StatusCode.INVALID_ARGUMENT, "Should return INVALID_ARGUMENT for empty display_name"

def test_get_zone_data_valid_zone(zone_data_service):
    request = ZoneDataRequest(zone_id="zone_abc")
    response = zone_data_service.GetZoneData(request)
    assert response.zone_id == "zone_abc"
    assert 0 <= response.coin_coeff <= 10, "Coin coeff out of expected range"
    assert len(response.display_name) > 0, "Display name should not be empty"

def test_get_zone_data_invalid_zone(zone_data_service):
    request = ZoneDataRequest(zone_id="")
    with pytest.raises(grpc.RpcError) as exc_info:
        zone_data_service.GetZoneData(request)
    assert exc_info.value.code() == grpc.StatusCode.INVALID_ARGUMENT, "Should return INVALID_ARGUMENT for empty zone_id"
