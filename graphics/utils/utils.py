import os
from os import path
import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns

metricsDir = "output"
graphicsDir = "plots"

methods = [
    "Chain",
    "Cuckoo",
    "Double",
    "Hopscotch",
    "Robin",
]

key_kinds = [
    "RandomKey",
    "SequentialKey"
]

load_factors = [
    "0.40",
    "0.60",
    "0.80",
]


def make_graphic(input_dir, output_dir, file_suffix, x, y, indexes):
    colors = sns.color_palette("husl", len(methods))
    plt.figure(figsize=(10, 8))

    for i, method in enumerate(methods):
        file_path = path.join(input_dir, method, f'{file_suffix}.csv')
        df = pd.read_csv(file_path, header=None)

        plt.plot(df[0], df[1], label=method, color=colors[i], linewidth=3)

    plt.xlabel(x)
    plt.xscale("log")
    plt.ylabel(y)
    plt.grid(True)
    plt.legend(loc='center left', bbox_to_anchor=(1, 0.5))
    plt.tight_layout()

    path_to_save = os.path.join(output_dir, f'{file_suffix}.png')
    os.makedirs(os.path.dirname(path_to_save), exist_ok=True)
    plt.savefig(path_to_save, dpi=1000)
    plt.close()


def make_two_graphic(input_dir, output_dir, file_suffix, x, y, indexes):
    colors = sns.color_palette("husl", len(methods))

    fig, (ax1, ax2) = plt.subplots(2, 1, figsize=(10, 8), sharex=True, height_ratios=[2, 2])

    m = ["Chain", "Double", "Hopscotch", "Robin"]
    for i, method in enumerate(m):
        file_path = path.join(input_dir, method, f'{file_suffix}.csv')
        df = pd.read_csv(file_path, header=None)

        ax1.plot(df[0], df[1], label=method, color=colors[i], linewidth=3)

    file_path = path.join(input_dir, "Cuckoo", f'{file_suffix}.csv')
    df = pd.read_csv(file_path, header=None)

    ax2.plot(df[indexes[0]], df[indexes[1]], label="Cuckoo", color=colors[len(colors) - 1], linewidth=3)

    ax1.set_xlabel(x)
    ax1.set_xscale("log")
    ax1.set_ylabel(y)
    ax1.legend(loc='center left', bbox_to_anchor=(1, 0.5))
    ax1.grid(True)

    ax2.set_xlabel(x)
    ax2.set_xscale("log")
    ax2.set_ylabel(y)
    ax2.legend(loc='center left', bbox_to_anchor=(1, 0.5))
    ax2.grid(True)

    plt.tight_layout()

    path_to_save = os.path.join(output_dir, f'{file_suffix}.png')
    os.makedirs(os.path.dirname(path_to_save), exist_ok=True)
    fig.savefig(path_to_save, dpi=1000)
    plt.close()
