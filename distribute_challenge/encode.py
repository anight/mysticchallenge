
from tblib import pickling_support
pickling_support.install()
import types
import pickle
import marshal

def request_serialize(func, depmodules, args, kwargs):
    # pickle is bad at functions serializing, let's use marshal
    code_serialized = marshal.dumps(func.__code__)
    return pickle.dumps( (code_serialized, depmodules, args, kwargs) )

def request_deserialize(binary, func_globals):
    code_serialized, depmodules, args, kwargs = pickle.loads(binary)
    code = marshal.loads(code_serialized)
    func = types.FunctionType(code, func_globals, "func")
    return func, depmodules, args, kwargs

def response_serialize(ret, stdout, stderr, exception):
    # marshal is bad at exceptions serializing, let's use pickle
    return pickle.dumps( (ret, stdout, stderr, exception) )

def response_deserialize(binary):
    ret, stdout, stderr, exception = pickle.loads(binary)
    return ret, stdout, stderr, exception
