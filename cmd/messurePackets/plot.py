import re
from collections import defaultdict
import matplotlib.pyplot as plt

log_files = {
    "Federated KLL Sketch": "kll-result",
    "Centralized KLL Sketch": "bad-kll-result",

    "Federated Count Sketch": "count-result",
    "Centralized Count Sketch": "bad-count-result",
}

header_pattern = re.compile(r"Running \w+ with MERGE_RATE=(\d+), CLIENT_AMOUNT=(\d+), STREAM_RATE=(\d+)")
first_number_pattern = re.compile(r"^(\d+)")

# final structured data exactly like old code
data = defaultdict(lambda: defaultdict(lambda: defaultdict(lambda: {
    "stream_rates": [],
    "processing_rates": []
})))

for dataset_name, file_path in log_files.items():
    try:
        with open(file_path, "r") as f:
            lines = f.readlines()
    except FileNotFoundError:
        print(f"Skipping missing file: {file_path}")
        continue

    current_merge = None
    current_client = None
    current_stream = None
    y_sum = 0
    collecting = False

    for line in lines:
        line = line.strip()

        header_match = header_pattern.search(line)
        if header_match:
            if collecting:
                # Save previous result
                processing_rate = y_sum * current_merge / 5
                total_stream = current_stream * current_client

                data[current_client][current_merge][dataset_name]["stream_rates"].append(total_stream)
                data[current_client][current_merge][dataset_name]["processing_rates"].append(processing_rate)

            current_merge = int(header_match.group(1))
            current_client = int(header_match.group(2))
            current_stream = int(header_match.group(3))
            y_sum = 0
            collecting = True
            continue

        num_match = first_number_pattern.match(line)
        if num_match and collecting:
            y_sum += int(num_match.group(1))

    # Save last run
    if collecting:
        processing_rate = y_sum * current_merge / 5
        total_stream = current_stream * current_client

        data[current_client][current_merge][dataset_name]["stream_rates"].append(total_stream)
        data[current_client][current_merge][dataset_name]["processing_rates"].append(processing_rate)

def plot_sketch_type(sketch_type, output_prefix):
    datasets = {
        f"Federated {sketch_type} Sketch": "Federated",
        f"Centralized {sketch_type} Sketch": "Centralized",
    }

    clients_values = [1, 10, 50, 100]
    merge_rates = [1000, 10000]

    colors = ['b', 'r']          # merge_rate → color
    linestyles = ['-', '--']     # federated → solid, centralized → dashed

    for client in clients_values:
        fig, ax = plt.subplots(figsize=(2.6, 2.6))

        handles_labels = []

        for ci, (dataset, label) in enumerate(datasets.items()):
            for mi, merge_rate in enumerate(merge_rates):

                sr = data[client][merge_rate][dataset]["stream_rates"]
                pr = data[client][merge_rate][dataset]["processing_rates"]

                if not sr:
                    continue

                # *** Multiply x-axis by 3 ***
                sr_scaled = [v * 3 for v in sr]

                line, = ax.plot(
                    sr_scaled,
                    pr,
                    linestyle=linestyles[ci],
                    marker='o',
                    markersize=2,
                    linewidth=1,
                    color=colors[mi],
                    label=f"{label} {sketch_type} (merge={merge_rate})"
                )
                handles_labels.append((line, line.get_label()))

        ax.set_title(f"Clients: {client*3}", fontsize=10)
        ax.set_xlabel("Stream Rate (T/s)", fontsize=9)
        ax.set_ylabel("Throughput (T/s)", fontsize=9)
        ax.grid(True)
        ax.ticklabel_format(style="sci", axis="both", scilimits=(0, 0))
        ax.tick_params(axis='both', which='both', labelsize=10)

        # *** Enforce y-axis = x-axis range ***
        max_x = ax.get_xlim()[1]
        ax.set_ylim(0, max_x)

        plt.tight_layout()

        # *** Save as individual plot ***
        filename = f"{output_prefix}-clients-{client}.png"
        plt.savefig(filename, bbox_inches="tight")
        plt.close()
        print(f"Saved: {filename}")



plot_sketch_type("KLL", "results-kll")
plot_sketch_type("Count", "results-count")
