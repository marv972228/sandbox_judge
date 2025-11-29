# Two Sum solution
# Reads array from line 1, target from line 2
# Outputs the two indices

def two_sum(nums, target):
    seen = {}
    for i, num in enumerate(nums):
        complement = target - num
        if complement in seen:
            return [seen[complement], i]
        seen[num] = i
    return []

# Read input
nums = list(map(int, input().split()))
target = int(input())

# Solve and print
result = two_sum(nums, target)
print(result[0], result[1])
