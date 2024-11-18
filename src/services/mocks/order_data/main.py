import grpc
from concurrent import futures
import time
import random
import os
import sys
import string
from grpc_reflection.v1alpha import reflection
import order_data_pb2
import order_data_pb2_grpc

def float_from_str(s):
    h = hash(s)
    max_hash = 2 ** (sys.hash_info.width - 1)
    return (h % max_hash) / float(max_hash)


def random_string_by_seed(seed, length):
    random.seed(seed)
    return ''.join(random.choice(string.ascii_lowercase) for _ in range(length))


class OrderDataServiceServicer(order_data_pb2_grpc.OrderDataServiceServicer):
    def GetOrderData(self, request, context):
        if len(request.order_id) == 0:
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            return order_data_pb2.OrderDataResponse()

        response = order_data_pb2.OrderDataResponse(
            order_id=request.order_id,
            user_id=random_string_by_seed(hash(request.order_id), 8),
            zone_id=random_string_by_seed(hash(request.order_id) + 1, 12),
            base_coin_amount=int(float_from_str(request.order_id) * 300),
        )
        return response


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    order_data_pb2_grpc.add_OrderDataServiceServicer_to_server(OrderDataServiceServicer(), server)
    SERVICE_NAMES = (
        order_data_pb2.DESCRIPTOR.services_by_name['OrderDataService'].full_name,
        reflection.SERVICE_NAME,
    )
    reflection.enable_server_reflection(SERVICE_NAMES, server)
    port = 50051
    server.add_insecure_port(f'[::]:{port}')
    server.start()
    print(f"Server started, listening on port {port}.")

    try:
        while True:
            time.sleep(86400)
    except KeyboardInterrupt:
        server.stop(0)

if __name__ == '__main__':
    serve()
