#!/usr/bin/env python3
"""
Naive O(n²) brute force solution - WILL TIMEOUT on large inputs!

This solution checks every pair of numbers.
For n=100,000, this is 10 billion operations.
Use this to verify the test cases actually catch slow solutions.
"""

def two_sum_naive(nums, target):
    n = len(nums)
    for i in range(n):
        for j in range(i + 1, n):
            if nums[i] + nums[j] == target:
                return [i, j]
    return []

# Read input
nums = list(map(int, input().split()))
target = int(input())

# Solve and print
result = two_sum_naive(nums, target)
print(result[0], result[1])
