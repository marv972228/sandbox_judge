#!/usr/bin/env python3
"""Wrong solution - always returns 0 0"""

def two_sum(nums, target):
    return [0, 0]  # Always wrong

if __name__ == "__main__":
    nums = list(map(int, input().split()))
    target = int(input())
    
    result = two_sum(nums, target)
    print(result[0], result[1])
