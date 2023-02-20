#! /usr/bin/env python3.11
import time

from distribute_challenge import compute_this

@compute_this(depmodules=[time]) # The Zen of Python says: Explicit is better than implicit
def func(x):
    time.sleep(x)
    return x*x

out = func(2).run()
assert out == 4

print("success")
