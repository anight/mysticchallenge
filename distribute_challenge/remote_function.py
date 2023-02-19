
import sys
import challenge_pb2
from .grpc import ServerContext
from .encode import request_serialize, response_deserialize

class MethodReturnedError(Exception): pass

class RemoteFunction:
    def __init__(self, func, server, depmodules):
        self.context = ServerContext(server)
        self.func = func
        self.depmodules = [ m.__name__ for m in depmodules ]

    def __call__(self, *args, **kwargs):
        return RemoteFunctionCall(self, *args, **kwargs)

class RemoteFunctionCall:
    def __init__(self, remote_func, *args, **kwargs):
        self.remote_func = remote_func
        self.args = args
        self.kwargs = kwargs

    def run(self):
        request_arg = request_serialize(
            self.remote_func.func, self.remote_func.depmodules, self.args, self.kwargs
        )

        request = self.remote_func.context.api.Execute(
            challenge_pb2.RequestExecute(
                request=request_arg
            )
        )

        response = self.remote_func.context.unary_unary_call(request)

        if response.error != "":
            raise MethodReturnedError(f"Execute() method failed: {response.error}")

        result, stdout, stderr, exception = response_deserialize(response.result)

        print(stdout, end='')
        print(stderr, end='', file=sys.stderr)

        if exception is not None:
            raise exception

        return result
