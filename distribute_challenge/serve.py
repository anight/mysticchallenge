
import io
import sys
import base64
from .encode import request_deserialize, response_serialize

def serve():
    print("ready", flush=True)
    while line := sys.stdin.readline():
        decoded_request = base64.b64decode(line)
        func, depmodules, args, kwargs = request_deserialize(decoded_request, globals())
        result, stdout, stderr, exception = care_wrapper(func, depmodules, *args, **kwargs)
        serialized_response = response_serialize(result, stdout, stderr, exception)
        encoded_response = base64.b64encode(serialized_response).decode('ascii')
        print(encoded_response, flush=True)

def care_wrapper(func, depmodules, *args, **kwargs):
    orig_stdout, orig_stderr = sys.stdout, sys.stderr
    sys.stdout = new_stdout = io.StringIO()
    sys.stderr = new_stderr = io.StringIO()
    result = None
    try:
        for m in depmodules:
            globals()[m] = __import__(m)
        result = func(*args, **kwargs)
    except Exception as e:
        sys.stdout = orig_stdout
        sys.stderr = orig_stderr
        return None, new_stdout.getvalue(), new_stderr.getvalue(), e
    sys.stdout = orig_stdout
    sys.stderr = orig_stderr
    return result, new_stdout.getvalue(), new_stderr.getvalue(), None

if __name__ == '__main__':
    serve()
