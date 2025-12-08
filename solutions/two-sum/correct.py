# Two Sum solution
# Reads array from line 1, target from line 2
# Outputs the two indices

def two_sum(nums, target):
    """
    Find two numbers that add up to target.
    
    THE PROBLEM:
    Given an array like [2, 7, 11, 15] and a target sum like 9,
    find which two numbers add up to 9. Return their positions (indices).
    Answer: 2 + 7 = 9, so return indices [0, 1]
    
    BRUTE FORCE APPROACH (slow):
    Check every pair: (2,7), (2,11), (2,15), (7,11), (7,15), (11,15)
    This takes O(n²) time - for 10,000 numbers, that's 100 million checks!
    
    OPTIMIZED APPROACH (what we use):
    Instead of searching for pairs, flip the question:
    "For each number, does the number I NEED already exist?"
    
    If target = 9 and current number = 2, then I NEED 7 (because 9 - 2 = 7)
    If target = 9 and current number = 7, then I NEED 2 (because 9 - 7 = 2)
    
    WHY USE A HASH MAP (dictionary)?
    A hash map lets us check "have I seen this number?" in O(1) time.
    - List search: "Is 7 in [2, 11, 15]?" → must check each element, O(n)
    - Hash map search: "Is 7 in {2: 0, 11: 2, 15: 3}?" → instant lookup, O(1)
    
    This is like the difference between:
    - Searching a phone book page by page (slow)
    - Looking up a name in a sorted index (fast)
    
    STEP-BY-STEP EXAMPLE:
    nums = [2, 7, 11, 15], target = 9
    
    Step 1: Look at 2 (index 0)
      - I need: 9 - 2 = 7
      - Have I seen 7? No (seen = {})
      - Remember 2's position: seen = {2: 0}
    
    Step 2: Look at 7 (index 1)
      - I need: 9 - 7 = 2
      - Have I seen 2? YES! At index 0
      - Found it! Return [0, 1]
    
    COMPLEXITY:
    - Time: O(n) - we look at each number at most once
    - Space: O(n) - we might store all numbers in the hash map
    """
    seen = {}  # Hash map: number -> index where we first saw it
    
    for i, num in enumerate(nums):
        complement = target - num  # What number do I need?
        
        if complement in seen:
            # Found it! The complement exists at seen[complement]
            # Current number is at index i
            # Together they sum to target
            return [seen[complement], i]
        
        # Haven't found a pair yet, remember this number's position
        # So if a future number needs THIS number, we can find it
        seen[num] = i
    
    return []  # No solution found (problem guarantees one exists)

# Read input
# Line 1: space-separated integers (the array)
# Line 2: the target sum
nums = list(map(int, input().split()))
target = int(input())

# Solve and print result as "index1 index2"
result = two_sum(nums, target)
print(result[0], result[1])
