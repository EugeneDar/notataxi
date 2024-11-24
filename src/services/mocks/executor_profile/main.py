import grpc
from concurrent import futures
import time
import random
import os
import sys
import string
from grpc_reflection.v1alpha import reflection
import executor_profile_pb2
import executor_profile_pb2_grpc

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


class ExecutorProfileServiceServicer(executor_profile_pb2_grpc.ExecutorProfileServiceServicer):
    def GetExecutorProfile(self, request, context):
        if len(request.display_name) == 0:
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            return executor_profile_pb2.ExecutorProfileResponse()

        response = executor_profile_pb2.ExecutorProfileResponse(
            id=str(hash(request.display_name) % 500000 + 1),
            tags=random_sublist(['fast', 'good conversation', 'good music', 'clear car'], 2, hash(request.display_name)),
            rating=float_from_str(request.display_name) + 4,
        )
        return response


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    executor_profile_pb2_grpc.add_ExecutorProfileServiceServicer_to_server(ExecutorProfileServiceServicer(), server)
    SERVICE_NAMES = (
        executor_profile_pb2.DESCRIPTOR.services_by_name['ExecutorProfileService'].full_name,
        reflection.SERVICE_NAME,
    )
    reflection.enable_server_reflection(SERVICE_NAMES, server)
    port = 9094
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
