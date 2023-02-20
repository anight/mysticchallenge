
import grpc

class ServerContext:
    def __init__(self, address, api_stub):
        self.address = address
        self.api = api_stub(self.new_channel())

    def new_channel(self):
        channel = None
        if '://' in self.address:
            scheme, hostport = self.address.split('://', 1)
            match scheme:
                case 'grpcs':
                    channel = grpc.aio.secure_channel(hostport, grpc.ssl_channel_credentials())
                case 'grpc':
                    channel = grpc.aio.insecure_channel(hostport)
        if channel is None:
            raise Exception("Server address must be in the form \"[grpc|grpcs]://host:port\"")
        return channel

    def unary_unary_call(self, call):
        async def future(call):
            await call.wait_for_connection()
            return await call
        return call._loop.run_until_complete(future(call))

