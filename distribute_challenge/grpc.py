
import grpc
import challenge_pb2_grpc

class ServerContext:
    def __init__(self, address):
        self.address = address
        self.api = self.new_service()

    def new_service(self):
        channel = None

        if '://' in self.address:
            proto, hostport = self.address.split('://', 1)
            match proto:
                case 'grpcs':
                    channel = grpc.aio.secure_channel(hostport, grpc.ssl_channel_credentials())
                case 'grpc':
                    channel = grpc.aio.insecure_channel(hostport)

        if channel is None:
            raise Exception("Server address must be in the form \"[grpc|grpcs]://host:port\"")

        return challenge_pb2_grpc.RemoteExecuteAPIStub(channel)

    def unary_unary_call(self, call):
        async def future(call):
            await call.wait_for_connection()
            return await call
        return call._loop.run_until_complete(future(call))

