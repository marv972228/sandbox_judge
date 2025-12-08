#!/usr/bin/env python3
"""
Generate test cases for two-sum problem.
Creates large inputs that will timeout naive O(n²) solutions.
"""

import random
import os

def generate_test(n, case_type="random"):
    """
    Generate a test case with n elements.
    
    case_type:
    - "random": random numbers, answer somewhere in middle
    - "worst_case": answer at the very end (maximizes work for naive solution)
    - "early_exit": answer at the beginning (best case)
    - "negative": includes negative numbers
    - "duplicates": has duplicate values
    """
    if case_type == "worst_case":
        # Answer is the last two elements - naive solution must check everything
        nums = list(range(1, n - 1))  # 1, 2, 3, ..., n-2
        # Last two numbers sum to a unique target
        nums.append(1000000)
        nums.append(2000000)
        target = 3000000
        answer = (n - 2, n - 1)
        
    elif case_type == "early_exit":
        # Answer is first two elements
        nums = [5, 10] + list(range(100, 100 + n - 2))
        target = 15
        answer = (0, 1)
        
    elif case_type == "negative":
        # Mix of positive and negative with UNIQUE answer
        # Use numbers that can't pair up to target except our answer
        nums = list(range(1, n - 1))  # All positive 1 to n-2
        # Add two large negative/positive that sum to unique target
        a, b = -999999, 1999999
        nums.append(a)
        nums.append(b)
        target = a + b  # 1000000 - unique, can't be made by small positives
        answer = (n - 2, n - 1)
        
    elif case_type == "duplicates":
        # Many duplicates but unique answer
        nums = [1] * (n // 2) + [2] * (n // 2 - 2) + [1000000, 2000000]
        target = 3000000
        answer = (n - 2, n - 1)
        
    else:  # random
        nums = [random.randint(1, 1000000) for _ in range(n - 2)]
        # Place answer at random positions
        idx1 = random.randint(0, n - 3)
        a = random.randint(1, 500000)
        b = random.randint(500001, 1000000)
        nums.insert(idx1, a)
        nums.append(b)
        target = a + b
        answer = (idx1, n - 1)
    
    return nums, target, answer

def write_test(path, nums, target, answer):
    """Write test input and expected output files."""
    with open(f"{path}.in", "w") as f:
        f.write(" ".join(map(str, nums)) + "\n")
        f.write(str(target) + "\n")
    
    with open(f"{path}.out", "w") as f:
        f.write(f"{answer[0]} {answer[1]}\n")

# Set seed for reproducibility
random.seed(42)

base_dir = "/home/marvi/repos/github.com/sandbox_judge/problems/two-sum/tests/hidden"

# Test 1: Medium size worst case (n=10,000)
# Naive O(n²) = 100 million operations, might be slow
nums, target, answer = generate_test(10000, "worst_case")
write_test(f"{base_dir}/1", nums, target, answer)
print(f"Test 1: n=10,000 worst case - answer at indices {answer}")

# Test 2: Large worst case (n=50,000)  
# Naive O(n²) = 2.5 billion operations, will definitely timeout
nums, target, answer = generate_test(50000, "worst_case")
write_test(f"{base_dir}/2", nums, target, answer)
print(f"Test 2: n=50,000 worst case - answer at indices {answer}")

# Test 3: Very large (n=100,000)
# Naive O(n²) = 10 billion operations, guaranteed timeout
nums, target, answer = generate_test(100000, "worst_case")
write_test(f"{base_dir}/3", nums, target, answer)
print(f"Test 3: n=100,000 worst case - answer at indices {answer}")

# Test 4: Negative numbers
nums, target, answer = generate_test(50000, "negative")
write_test(f"{base_dir}/4", nums, target, answer)
print(f"Test 4: n=50,000 with negatives - target={target}")

# Test 5: Many duplicates
nums, target, answer = generate_test(50000, "duplicates")
write_test(f"{base_dir}/5", nums, target, answer)
print(f"Test 5: n=50,000 with duplicates")

# Test 6: Edge case - minimum size
nums, target, answer = [3, 5], 8, (0, 1)
write_test(f"{base_dir}/6", nums, target, answer)
print(f"Test 6: n=2 minimum size")

print("\nAll hidden test cases generated!")
