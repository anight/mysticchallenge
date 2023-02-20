#! /usr/bin/env python3.11
import os
import time
import random
import asyncio
import collections

loop = asyncio.new_event_loop()
asyncio.set_event_loop(loop)

from distribute_challenge import compute_this

@compute_this(depmodules=[os, time]) # The Zen of Python says: Explicit is better than implicit
def func(x):
    time.sleep(x)
    return x*x, os.getpid()

server_workers = func.get_workers()

def test(factor):
    calls = server_workers * factor
    numbers = [ random.random() for _ in range(calls) ]

    async def load_test():
        result = await asyncio.gather(*[ func(x).async_run() for x in numbers ])
        squares, pids = zip(*result)
        return sum(squares), pids

    started = time.time()
    sum_of_squares, pids = loop.run_until_complete(load_test())
    elapsed = time.time() - started

    assert sum_of_squares == sum(x*x for x in numbers)

    minimal_execution_time = max(max(numbers), sum(numbers) / server_workers)

    print("-" * 80)
    print(f"workers: {server_workers}")
    print(f"calls: {calls}")
    print(f"minimal_execution_time: {minimal_execution_time:.6f}")
    print(f"real_execution_time: {elapsed:.6f}")
    print(f"effectiveness of parallelization: {100 * minimal_execution_time / elapsed:.2f}%")

    pid_stats = dict(collections.Counter(pids)).values()
    hist = dict(collections.Counter(pid_stats))

    print()
    print("num of calls: served by num of wokers")
    for k, v in sorted(hist.items()):
        print(f"{k:12d}: {v:4d}")

for factor in (1, 2, 5, 10):
    test(factor)
