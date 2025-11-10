import re
import matplotlib.pyplot as plt
from collections import defaultdict

# Path to the log file
file_path = "kll-result"

# Regex patterns
header_pattern = re.compile(r"Running \w+ with MERGE_RATE=(\d+), CLIENT_AMOUNT=(\d+), STREAM_RATE=(\d+)")
first_number_pattern = re.compile(r"^(\d+)")  # First number in the numeric row

# Data structure: data[client_amount][merge_rate] = list of (stream_rate, y_sum)
data = defaultdict(lambda: defaultdict(list))

with open(file_path, "r") as f:
    lines = f.readlines()

current_merge = None
current_client = None
current_stream = None
y_sum = 0
collecting = False

for line in lines:
    line = line.strip()
    
    # Check for a header
    header_match = header_pattern.search(line)
    if header_match:
        # Save previous run if there was one
        if collecting and current_client is not None:
            data[current_client][current_merge].append((current_stream, y_sum))
        
        # Start new run
        current_merge = int(header_match.group(1))
        current_client = int(header_match.group(2))
        current_stream = int(header_match.group(3))
        y_sum = 0
        collecting = True
        continue
    
    # Check for numeric row (first number only)
    number_match = first_number_pattern.match(line)
    if number_match and collecting:
        y_sum += int(number_match.group(1))
        
# Save the last run
if collecting and current_client is not None:
    data[current_client][current_merge].append((current_stream, y_sum))

# Plotting
for client_amount, merge_dict in data.items():
    plt.figure(figsize=(8, 5))
    for merge_rate, values in merge_dict.items():
        # Sort by stream_rate
        values.sort(key=lambda x: x[0])
        x = [v[0] for v in values]
        y = [v[1]*merge_rate for v in values]
        plt.plot(x, y, marker='o', label=f'MergeRate={merge_rate}')

    plt.title(f'ClientAmount = {client_amount}')
    plt.xlabel('StreamRate')
    plt.ylabel('Sum of first numbers across all nodes (y)')
    plt.legend()
    plt.grid(True)
    plt.tight_layout()
    plt.show()
