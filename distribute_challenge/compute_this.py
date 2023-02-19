
from .remote_function import RemoteFunction

def compute_this(server='grpc://127.0.0.1:8000', depmodules=[]):
    def decorator(func):
        return RemoteFunction(func, server, depmodules)
    return decorator
