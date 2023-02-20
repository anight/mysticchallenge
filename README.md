# Software components choices
For the server side language I picked Golang. It is well known for its simplicity, scalability, built-in parallelization means, performance and relatively inexpensive engineers available on the market. It also suits well for the purpose of the challenge.

GRPC was chosen over all other transports because it is an industry standard for RPC nowdays. A typical Go/grpc service can easily handle tens of thousands concurrent tcp connections without noticable servic degradation. The whole grpc infrastructure makes it a great choice for rapid development of projects of any size. The performance of GRPC is acceptable for most common uses. Although 13 years ago I have developed my own RPC system based on protobuf which was up to x3 times faster on some workload compared to GRPC.

GRPC bindings for python have support for asyncio interface since version 1.32. Asyncio interface is much easier to use when a request parallelization is required, like in the challenge.

# Build and setup
To build the server:
```bash
$ go get && go build
```
To install needed python modules for distribute_challenge package:
```bash
$ pip3 install --user -r ./distribute_challenge/requirements.txt
```

# Server start
```bash
$ ./mysticchallenge -workers 20
```

# Tests start
```bash
$ ./test.py
success
$ ./test_load.py
--------------------------------------------------------------------------------
workers: 20
calls: 20
minimal_execution_time: 0.999584
real_execution_time: 1.001084
effectiveness of parallelization: 99.85%

num of calls: served by num of wokers
           1:   20
--------------------------------------------------------------------------------
workers: 20
calls: 40
minimal_execution_time: 1.022768
real_execution_time: 1.598288
effectiveness of parallelization: 63.99%

num of calls: served by num of wokers
           1:    5
           2:   10
           3:    5
--------------------------------------------------------------------------------
workers: 20
calls: 100
minimal_execution_time: 2.588694
real_execution_time: 3.130264
effectiveness of parallelization: 82.70%

num of calls: served by num of wokers
           3:    2
           4:    7
           5:    5
           6:    1
           7:    5
--------------------------------------------------------------------------------
workers: 20
calls: 200
minimal_execution_time: 5.141060
real_execution_time: 5.615868
effectiveness of parallelization: 91.55%

num of calls: served by num of wokers
           7:    2
           8:    3
           9:    2
          10:    4
          11:    6
          12:    1
          13:    2
```

# Server: design and implementation

At the start server creates a number of python worker processes. In case of worker processes crash it restarts them automatically. Request scheduling to workers is implemented with a help of channels. From writer side of the channel there are goroutines serving grpc requests. From reader side of the channel there are `worker.serve()` goroutines. All reads and writes are guaranteed to be atomic, so it does not matter if goroutines are physically executed in different OS threads. If there is at least one free worker reading from the queue then the job will be scheduled within microseconds. Otherwise the job will wait either in the queue channel or (in case if the queue is full) it will block queue channel write operation from grpc request goroutine side.
Server does not know anything specific about python. It operates only with opaque requests and responses to and from workers. When server takes a request for remote function execution it does the following:
 - base64-encoded request is being sent as one line into stdin of the worker
 - as soon as response is ready, base64-encoded response is sent back as one line from stdout of the worker

Server does have some error condition checks although for production system there should be much more, like proper request timeout managment and so on.
All stderr from worker processes are logged into server logs.

# Client: design and implementation

From python point of view every request to the service consists of:
 - function
 - depmodules (a list of modules which needed to be imported before function execution)
 - *args
 - **kwargs

Every response consists of:
 - returned value from the function
 - stdout (string generated by the function)
 - stderr (string generated by the function)
 - exception generated by the function

Please note that stderr output of the remote function is not the same as stderr of the worker process. The former is captured during function execution and sent back with response and the latter is any unexpected but still potentially usefull debug message that are stored in server logs.

Requests and responeses are serialized with a help of `pickle` and `marshal`. The functionality of the two modules is pretty much the same, however pickle has no means to serialize a function for the challenge purposes (or at least I could not find a proper way to make it work). On the other side `marshal` is unable to serialize exceptions, so I had to use `pickle` together with `tblib`.

I added `func.async_run()` method to support asynchronous calls.

I also added `GetWorkers` method to GRPC protocol in order to simplify tests.

# test_load.py
The purpose of the test is to demonstrate how requests are distributed between workers. Multiple tests are done for different number of function calls. For different number of workers one should restart server with different `-workers` value.

In `test_load_stats.txt` I provide results of multiple test_load.py runs for 10, 20, 50, 100 and 200 workers.
For every run there is a histogram showing how many workers executed how many requests.
`minimal_execution_time` is the total time spent by `time.sleep()` function divided by number of available workers but not less than the longest time of all calls.
`effectiveness of parallelization` is percentage of `minimal_execution_time` to total time of the test.

# Finale
I think there are a lot of bugs in my code, but at least it works within conditions I mentioned.

 - 166 lines of python code in distribute_challenge module
 - 293 lines of go code for server

Your feedback is very welcome !