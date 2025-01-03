import grpc
from concurrent import futures
import time
import random
import os
import sys
from grpc_reflection.v1alpha import reflection
import zone_data_pb2
import zone_data_pb2_grpc

class ZoneDataServiceServicer(zone_data_pb2_grpc.ZoneDataServiceServicer):
    def GetZoneData(self, request, context):
        if len(request.zone_id) == 0:
            context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
            return zone_data_pb2.ZoneDataResponse()

        params = random.choice([
            ('Lyubertsy', 1.5), 
            ('Severodvinsk', 1.0), 
            ('Barysaw', 1.0), 
            ('Moscow', 2.0), 
            ('Kamensk-Uralsky', 1.0)]
        )
        
        response = zone_data_pb2.ZoneDataResponse(
            zone_id=request.zone_id,
            coin_coeff=params[1],
            display_name=params[0]
        )
        return response

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    zone_data_pb2_grpc.add_ZoneDataServiceServicer_to_server(ZoneDataServiceServicer(), server)
    SERVICE_NAMES = (
        zone_data_pb2.DESCRIPTOR.services_by_name['ZoneDataService'].full_name,
        reflection.SERVICE_NAME,
    )
    reflection.enable_server_reflection(SERVICE_NAMES, server)
    port = 9092
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
