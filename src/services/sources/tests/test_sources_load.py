import os
import sys
from locust import User, task, between, events
import grpc
import time
import random
from datetime import datetime

sys.path.append(os.path.join(os.path.dirname(os.path.abspath(__file__)), '../protobufs'))

from sources_pb2_grpc import SourcesServiceStub
from sources_pb2 import SourcesRequest

class GRPCUser(User):
    abstract = True
    wait_time = between(0.1, 1)

    def __init__(self, environment):
        super().__init__(environment)
        self.client = None
        
    def on_start(self):
        self.channel = grpc.insecure_channel("130.193.46.187:9000")
        self.client = SourcesServiceStub(self.channel)
        self.order_ids = [f"order_{i}" for i in range(1000)]
        self.executor_ids = [f"exec_{i}" for i in range(100)]

    def on_stop(self):
        if self.channel:
            self.channel.close()

class SourcesUser(GRPCUser):
    @task
    def get_order_info(self):
        order_id = random.choice(self.order_ids)
        executor_id = random.choice(self.executor_ids)
        request = SourcesRequest(
            order_id=order_id,
            executor_id=executor_id
        )

        start_time = time.time()
        try:
            self.client.GetOrderInfo(request, timeout=5)
            events.request.fire(
                request_type="grpc",
                name="GetOrderInfo",
                response_time=(time.time() - start_time) * 1000,
                response_length=0,
                exception=None,
                context=None
            )
        except grpc.RpcError as e:
            events.request.fire(
                request_type="grpc",
                name="GetOrderInfo",
                response_time=(time.time() - start_time) * 1000,
                response_length=0,
                exception=e,
                context=None
            )