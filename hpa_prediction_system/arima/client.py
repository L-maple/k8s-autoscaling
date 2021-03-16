import grpc
import forecast_pb2_grpc
import forecast_pb2


def run():
    with grpc.insecure_channel("localhost:50000") as channel:
        stub = forecast_pb2_grpc.ForecastServiceStub(channel)
        request = forecast_pb2.ForecastRequest(
            data="hello",
            minutes=1,
        )
        response = stub.GetForeCastValue(request)
        print(response.value)


if __name__ == "__main__":
    run()
