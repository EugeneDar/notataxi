import grpc
from concurrent import futures
import time
import random
import os
import sys
import string
from grpc_reflection.v1alpha import reflection
sys.path.append(os.path.join(os.path.dirname(os.path.abspath(__file__)), '../../sources/protobufs'))
import config_pb2
import config_pb2_grpc

def float_from_str(s):
    h = hash(s)
    max_hash = 2 ** (sys.hash_info.width - 1)
    return (h % max_hash) / float(max_hash)


def random_string_by_seed(seed, length):
    random.seed(seed)
    return ''.join(random.choice(string.ascii_lowercase) for _ in range(length))


def random_sublist(list, subset_size, seed):
    random.seed(seed)
    return random.sample(list, subset_size)


class ConfigServiceServicer(config_pb2_grpc.ConfigServiceServicer):
    def GetConfig(self, request, context):
        response = config_pb2.ConfigResponse(
            settings={
                'coin_coeff_settings_maximum': '3',
                'coin_coeff_settings_fallback': '1',
            },
        )
        return response


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    config_pb2_grpc.add_ConfigServiceServicer_to_server(ConfigServiceServicer(), server)
    SERVICE_NAMES = (
        config_pb2.DESCRIPTOR.services_by_name['ConfigService'].full_name,
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
