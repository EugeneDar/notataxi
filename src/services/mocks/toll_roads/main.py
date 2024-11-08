import grpc
from concurrent import futures
import time
import random
import os
import sys
from grpc_reflection.v1alpha import reflection
sys.path.append(os.path.join(os.path.dirname(os.path.abspath(__file__)), '../../sources/protobufs'))
import toll_roads_pb2
import toll_roads_pb2_grpc

def float_from_str(s):
    h = hash(s)
    max_hash = 2 ** (sys.hash_info.width - 1)
    return (h % max_hash) / float(max_hash)


class TollRoadsServiceServicer(toll_roads_pb2_grpc.TollRoadsServiceServicer):
    def GetTollRoads(self, request, context):
        if len(request.display_name) == 0:
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            return toll_roads_pb2.TollRoadsResponse()

        response = toll_roads_pb2.TollRoadsResponse(
            bonus_amount=int(float_from_str(request.display_name) * 100)
        )
        return response


def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    toll_roads_pb2_grpc.add_TollRoadsServiceServicer_to_server(TollRoadsServiceServicer(), server)
    SERVICE_NAMES = (
        toll_roads_pb2.DESCRIPTOR.services_by_name['TollRoadsService'].full_name,
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
