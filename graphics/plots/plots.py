import os
import pandas as pd
import seaborn as sns
from os import path
from typing import Tuple, Optional, List
from matplotlib import pyplot as plt

# Configuration constants
METHODS = ["Chain", "Cuckoo", "Double", "Hopscotch", "Robin"]
METHODS_FULL_NAME = {
    "Chain": "Chain method",
    "Cuckoo": "Cuckoo method",
    "Double": "Double hashing",
    "Hopscotch": "Hopscotch hashing",
    "Robin": "Robin method hashing"
}
KEY_KINDS = ["RandomKey", "SequentialKey"]
LOAD_FACTORS = ["0.40", "0.60", "0.80"]


def make_plot(
        input_dir: str,
        output_dir: str,
        x_label: str,
        y_label: str,
        data_indexes: Tuple[int, int],
        split_methods: Optional[List[str]] = None,
        use_log_scale_x: bool = False,
        use_log_scale_y: bool = False
) -> None:
    """
    Function to create plots for different scenarios.
    
    Args:
        input_dir: Directory containing input CSV files
        output_dir: Directory to save the output plot
        x_label: Label for x-axis
        y_label: Label for y-axis
        data_indexes: Tuple of (x_index, y_index) for data columns
        split_methods: Optional list of methods to plot in a separate subplot
        use_log_scale_x: Whether to use logarithmic scale for x-axis
        use_log_scale_y: Whether to use logarithmic scale for y-axis
    """
    for key_kind in KEY_KINDS:
        for load_factor in LOAD_FACTORS:
            file_suffix = f'{key_kind}_{load_factor}'

            if load_factor == "0.40" or not split_methods:
                make_graphic(
                    input_dir=input_dir,
                    output_dir=output_dir,
                    file_suffix=file_suffix,
                    x_label=x_label,
                    y_label=y_label,
                    data_indexes=data_indexes,
                    methods=METHODS,
                    use_log_scale_x=use_log_scale_x,
                    use_log_scale_y=use_log_scale_y
                )
            else:
                make_graphic(
                    input_dir=input_dir,
                    output_dir=output_dir,
                    file_suffix=file_suffix,
                    x_label=x_label,
                    y_label=y_label,
                    data_indexes=data_indexes,
                    methods=METHODS,
                    split_methods=split_methods,
                    use_log_scale_x=use_log_scale_x,
                    use_log_scale_y=use_log_scale_y
                )


def make_insert_no_reserve_time_graphics():
    make_plot(
        input_dir=path.join("data", "InsertNoReserve"),
        output_dir=path.join("graphics", "InsertNoReserve"),
        x_label="Количество элементов",
        y_label="Среднее время 1 операции (ns)",
        data_indexes=(0, 1),
        split_methods=["Cuckoo"],
        use_log_scale_x=True,
    )


def make_insert_reserve_time_graphics():
    make_plot(
        input_dir=path.join("data", "InsertReserve"),
        output_dir=path.join("graphics", "InsertReserve"),
        x_label="Количество элементов",
        y_label="Среднее время 1 операции (ns)",
        data_indexes=(0, 1),
        use_log_scale_x=True,
    )


def make_insert_reserve_memory_graphics():
    make_plot(
        input_dir=path.join("data", "InsertNoReserve"),
        output_dir=path.join("graphics", "AllocateMemoryInsertNoReserve"),
        x_label="Количество элементов",
        y_label="Количество выделенной памяти (bytes)",
        data_indexes=(0, 2),
        split_methods=["Cuckoo"],
        use_log_scale_x=True,
    )


def make_success_get_graphics():
    make_plot(
        input_dir=path.join("data", "SuccessGet"),
        output_dir=path.join("graphics", "SuccessGet"),
        x_label="Количество элементов",
        y_label="Среднее время 1 операции (ns)",
        data_indexes=(0, 1),
        use_log_scale_x=True,
    )


def make_unsuccess_get_graphics():
    make_plot(
        input_dir=path.join("data", "UnsuccessGet"),
        output_dir=path.join("graphics", "UnsuccessGet"),
        x_label="Количество элементов",
        y_label="Среднее время 1 операции (ns)",
        data_indexes=(0, 1),
        use_log_scale_x=True,
    )


def make_delete_graphics():
    make_plot(
        input_dir=path.join("data", "Delete"),
        output_dir=path.join("graphics", "Delete"),
        x_label="Количество элементов",
        y_label="Среднее время 1 операции (ns)",
        data_indexes=(0, 1),
        use_log_scale_x=True,
    )


def make_graphic(
        input_dir: str,
        output_dir: str,
        file_suffix: str,
        x_label: str,
        y_label: str,
        data_indexes: Tuple[int, int],
        methods: List[str],
        split_methods: Optional[List[str]] = None,
        height_ratios: Optional[List[int]] = None,
        dpi: int = 1000,
        figsize: Tuple[int, int] = (10, 8),
        use_log_scale_x: bool = False,
        use_log_scale_y: bool = False
) -> None:
    """
    Create a single or split plot based on the provided configuration.
    
    Args:
        input_dir: Directory containing input CSV files
        output_dir: Directory to save the output plot
        file_suffix: Suffix for the input/output files
        x_label: Label for x-axis
        y_label: Label for y-axis
        data_indexes: Tuple of (x_index, y_index) for data columns
        methods: List of methods to plot
        split_methods: Optional list of methods to plot in a separate subplot
        height_ratios: Optional list of height ratios for subplots
        dpi: DPI for the output image
        figsize: Figure size as (width, height)
        use_log_scale_x: Whether to use logarithmic scale for x-axis
        use_log_scale_y: Whether to use logarithmic scale for y-axis
    """
    colors = sns.color_palette("husl", len(methods))

    if split_methods:
        # Create split plot
        fig, (ax1, ax2) = plt.subplots(2, 1, figsize=figsize, sharex=True,
                                       height_ratios=height_ratios or [2, 2])

        # Plot main methods
        for i, method in enumerate(methods):
            if method not in split_methods:
                file_path = path.join(input_dir, method, f'{file_suffix}.csv')
                df = pd.read_csv(file_path, header=None)
                ax1.plot(df[data_indexes[0]], df[data_indexes[1]],
                         label=METHODS_FULL_NAME[method], color=colors[i], linewidth=3)

        # Plot split methods
        for method in split_methods:
            file_path = path.join(input_dir, method, f'{file_suffix}.csv')
            df = pd.read_csv(file_path, header=None)
            ax2.plot(df[data_indexes[0]], df[data_indexes[1]],
                     label=METHODS_FULL_NAME[method], color=colors[methods.index(method)], linewidth=3)

        # Configure subplots
        for ax in [ax1, ax2]:
            ax.set_xlabel(x_label, fontsize=18)
            if use_log_scale_x:
                ax.set_xscale("log")
            ax.set_ylabel(y_label, fontsize=18)
            if use_log_scale_y:
                ax.set_yscale("log")
            ax.tick_params(axis='both', labelsize=16)
            ax.yaxis.offsetText.set_fontsize(18)
            ax.legend(loc='best', fontsize=17, frameon=True, borderpad=1.2)
            ax.grid(True)
    else:
        # Create single plot
        plt.figure(figsize=figsize)

        for i, method in enumerate(methods):
            file_path = path.join(input_dir, method, f'{file_suffix}.csv')
            df = pd.read_csv(file_path, header=None)
            plt.plot(df[data_indexes[0]], df[data_indexes[1]],
                     label=METHODS_FULL_NAME[method], color=colors[i], linewidth=3)

        plt.xlabel(x_label, fontsize=18)
        if use_log_scale_x:
            plt.xscale("log")
        plt.ylabel(y_label, fontsize=18)
        if use_log_scale_y:
            plt.yscale("log")
        plt.grid(True)
        plt.tick_params(axis='both', labelsize=16)
        plt.legend(loc='best', fontsize=17, frameon=True, borderpad=1.2)

    plt.tight_layout()

    # Save plot
    path_to_save = os.path.join(output_dir, f'{file_suffix}.png')
    os.makedirs(os.path.dirname(path_to_save), exist_ok=True)
    plt.savefig(path_to_save, dpi=dpi)
    plt.close()
