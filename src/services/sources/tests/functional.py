import pytest
import grpc
from google.protobuf import empty_pb2

from src.services.sources.protobufs.config_pb2_grpc import ConfigServiceStub
from src.services.sources.protobufs.config_pb2 import ConfigResponse
from src.services.sources.protobufs.executor_profile_pb2_grpc import ExecutorProfileServiceStub
from src.services.sources.protobufs.executor_profile_pb2 import ExecutorProfileRequest
from src.services.sources.protobufs.order_data_pb2_grpc import OrderDataServiceStub
from src.services.sources.protobufs.order_data_pb2 import OrderDataRequest
from src.services.sources.protobufs.sources_pb2_grpc import OrderInfoServiceStub
from src.services.sources.protobufs.sources_pb2 import OrderInfoRequest
from src.services.sources.protobufs.toll_roads_pb2_grpc import TollRoadsServiceStub
from src.services.sources.protobufs.toll_roads_pb2 import TollRoadsRequest
from src.services.sources.protobufs.zone_data_pb2_grpc import ZoneDataServiceStub
from src.services.sources.protobufs.zone_data_pb2 import ZoneDataRequest

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

def test_get_config(config_service):
    response = config_service.GetConfig(empty_pb2.Empty())
    # assert isinstance(response, ConfigResponse)
    assert isinstance(response.coin_coeff_settings.maximum, float)
    assert isinstance(response.coin_coeff_settings.fallback, float)
    assert 0 <= response.coin_coeff_settings.maximum <= 10, "Maximum coeff out of range"
    assert 0 <= response.coin_coeff_settings.fallback <= 10, "Fallback coeff out of range"

def test_get_executor_profile_valid_request(executor_profile_service):
    request = ExecutorProfileRequest(display_name="Valid User")
    response = executor_profile_service.GetExecutorProfile(request)
    assert response.id
    assert len(response.id) > 0, "ID should not be empty"
    assert isinstance(response.tags, list)
    assert all(isinstance(tag, str) for tag in response.tags), "All tags should be strings"
    assert 0 <= response.rating <= 5, "Rating out of range"

def test_get_executor_profile_empty_name(executor_profile_service):
    request = ExecutorProfileRequest(display_name="")
    response = executor_profile_service.GetExecutorProfile(request)
    assert response.id == "", "ID should be empty for unknown display_name"
    assert response.rating == 0.0, "Rating should be zero for unknown display_name"
    assert not response.tags, "Tags should be empty for unknown display_name"

def test_get_order_data_valid_order(order_data_service):
    request = OrderDataRequest(order_id="order_123")
    response = order_data_service.GetOrderData(request)
    assert response.order_id == "order_123"
    assert len(response.user_id) > 0, "User ID should not be empty"
    assert len(response.zone_id) > 0, "Zone ID should not be empty"
    assert response.base_coin_amount > 0, "Base coin amount should be positive"

def test_get_order_data_invalid_order(order_data_service):
    request = OrderDataRequest(order_id="invalid_order")
    with pytest.raises(grpc.RpcError) as exc_info:
        order_data_service.GetOrderData(request)
    assert exc_info.value.code() == grpc.StatusCode.NOT_FOUND, "Should return NOT_FOUND for invalid order"

def test_get_order_info_valid_order(order_info_service):
    request = OrderInfoRequest(order_id="order_123", executor_id="exec_456")
    response = order_info_service.GetOrderInfo(request)
    assert response.order_id == "order_123"
    assert response.final_coin_amount > 0, "Final coin amount should be positive"
    assert response.price_components.base_coin_amount > 0, "Base coin amount should be positive"
    assert response.price_components.coin_coeff >= 1, "Coin coefficient should be >= 1"
    assert isinstance(response.executor_profile.rating, float)
    assert 0 <= response.executor_profile.rating <= 5, "Rating out of range"

def test_get_order_info_invalid_executor(order_info_service):
    request = OrderInfoRequest(order_id="order_123", executor_id="invalid_exec")
    with pytest.raises(grpc.RpcError) as exc_info:
        order_info_service.GetOrderInfo(request)
    assert exc_info.value.code() == grpc.StatusCode.NOT_FOUND, "Should return NOT_FOUND for invalid executor_id"

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

def test_get_zone_data_nonexistent_zone(zone_data_service):
    request = ZoneDataRequest(zone_id="nonexistent_zone")
    with pytest.raises(grpc.RpcError) as exc_info:
        zone_data_service.GetZoneData(request)
    assert exc_info.value.code() == grpc.StatusCode.NOT_FOUND, "Should return NOT_FOUND for nonexistent zone"