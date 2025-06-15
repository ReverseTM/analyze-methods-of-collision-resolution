import csv
import os
import re
from os import path
from typing import Pattern

# Configuration constants
METHODS = ["Chain", "Cuckoo", "Double", "Hopscotch", "Robin"]
KEY_KINDS = ["RandomKey", "SequentialKey"]
LOAD_FACTORS = ["0.40", "0.60", "0.80"]


def parse_results() -> None:
    """
    Parse all benchmark results and save them to CSV files.
    
    This function iterates through all combinations of methods, key kinds, and load factors,
    and parses the following types of data:
    - Insert operations (with and without reserve)
    - Get operations (unsuccessful)
    - Delete operations
    - Collisions data
    
    The results are saved in the 'data' directory with the following structure:
    data/
    ├── InsertNoReserve/
    │   └── {method}/
    │       └── {key_kind}_{load_factor}.csv
    ├── InsertReserve/
    │   └── {method}/
    │       └── {key_kind}_{load_factor}.csv
    ├── SuccessGet/
    │   └── {method}/
    │       └── {key_kind}_{load_factor}.csv
    ├── UnsuccessGet/
    │   └── {method}/
    │       └── {key_kind}_{load_factor}.csv
    ├── Delete/
    │   └── {method}/
    │       └── {key_kind}_{load_factor}.csv
    └── Collisions/
        └── {method}/
            └── {key_kind}_{load_factor}.csv
    """
    for method in METHODS:
        for key_kind in KEY_KINDS:
            for load_factor in LOAD_FACTORS:
                parse_insert_to_csv(method=method, key_kind=key_kind, load_factor=load_factor)
                parse_get_to_csv(method=method, key_kind=key_kind, load_factor=load_factor)
                parse_delete_to_csv(method=method, key_kind=key_kind, load_factor=load_factor)
                parse_collisions(method=method, key_kind=key_kind, load_factor=load_factor)


def parse_insert_to_csv(method: str, key_kind: str, load_factor: str) -> None:
    """
    Parse insert benchmark results and save to CSV.
    
    Args:
        method: Hash table implementation method
        key_kind: Type of keys used in benchmark
        load_factor: Load factor used in benchmark
    """
    for operation in ["InsertNoReserve", "InsertReserve"]:
        input_file_path = path.join("results", operation, method, "row.txt")
        output_file_path = path.join("data", operation, method, f'{key_kind}_{load_factor}.csv')

        pattern = re.compile(
            rf"Benchmark{re.escape(operation)}/{re.escape(method)}-{re.escape(key_kind)}-{load_factor}-(?P<Size>\d+)-\d+\s+"
            r"\d+\s+"
            r"[\d.]+\s+ns/op\s+"
            r"(?P<NsPerInsert>[\d.]+)\s+ns/insert\s+"
            r"(?P<BytesPerOp>\d+)\s+B/op\s+"
            r"\d+\s+allocs/op"
        )

        def extract_row(match):
            return [
                int(match.group('Size')),
                float(match.group('NsPerInsert')),
                int(match.group('BytesPerOp'))
            ]

        parse_benchmark_results(input_file_path, output_file_path, pattern, extract_row)


def parse_get_to_csv(method: str, key_kind: str, load_factor: str) -> None:
    """
    Parse get benchmark results and save to CSV.
    
    Args:
        method: Hash table implementation method
        key_kind: Type of keys used in benchmark
        load_factor: Load factor used in benchmark
    """
    for operation in ["SuccessGet", "UnsuccessGet"]:
        input_file_path = path.join("results", operation, method, "row.txt")
        output_file_path = path.join("data", operation, method, f'{key_kind}_{load_factor}.csv')

        pattern = re.compile(
            rf"Benchmark{re.escape(operation)}/{re.escape(method)}-{re.escape(key_kind)}-{load_factor}-(?P<Size>\d+)-\d+\s+"
            r"\d+\s+"
            r"(?P<NsPerOp>[\d.]+)\s+ns/op\s+"
        )

        def extract_row(match):
            return [
                int(match.group('Size')),
                float(match.group('NsPerOp'))
            ]

        parse_benchmark_results(input_file_path, output_file_path, pattern, extract_row)


def parse_delete_to_csv(method: str, key_kind: str, load_factor: str) -> None:
    """
    Parse delete benchmark results and save to CSV.
    
    Args:
        method: Hash table implementation method
        key_kind: Type of keys used in benchmark
        load_factor: Load factor used in benchmark
    """
    operation = "Delete"
    input_file_path = path.join("results", operation, method, "row.txt")
    output_file_path = path.join("data", operation, method, f'{key_kind}_{load_factor}.csv')

    pattern = re.compile(
        rf"Benchmark{re.escape(operation)}/{re.escape(method)}-{re.escape(key_kind)}-{load_factor}-(?P<Size>\d+)-\d+\s+"
        r"\d+\s+"
        r"(?P<NsPerOp>[\d.]+)\s+ns/op\s+"
    )

    def extract_row(match):
        return [
            int(match.group('Size')),
            float(match.group('NsPerOp'))
        ]

    parse_benchmark_results(input_file_path, output_file_path, pattern, extract_row)


def parse_collisions(method: str, key_kind: str, load_factor: str) -> None:
    """
    Parse collisions data and save to CSV.
    
    Args:
        method: Hash table implementation method
        key_kind: Type of keys used in benchmark
        load_factor: Load factor used in benchmark
    """
    operation = "Collisions"
    input_file_path = path.join("results", operation, method, load_factor, f'{key_kind}.csv')
    output_file_path = path.join("data", operation, method, f'{key_kind}_{load_factor}.csv')

    copy_csv_file(input_file_path, output_file_path)


def copy_csv_file(input_file_path: str, output_file_path: str) -> None:
    """
    Copy CSV file from input to output location.
    
    Args:
        input_file_path: Path to input file
        output_file_path: Path to output file
    """
    ensure_output_dir(output_file_path)
    
    with open(input_file_path, "r", newline='') as input_file, \
         open(output_file_path, 'w', newline='') as output_file:
        reader = csv.reader(input_file)
        writer = csv.writer(output_file)
        writer.writerows(reader) 


def ensure_output_dir(output_file_path: str) -> None:
    """Create output directory if it doesn't exist."""
    os.makedirs(os.path.dirname(output_file_path), exist_ok=True)


def parse_benchmark_results(
    input_file_path: str,
    output_file_path: str,
    pattern: Pattern,
    extract_row: callable
) -> None:
    """
    Parse benchmark results using regex pattern and save to CSV.
    
    Args:
        input_file_path: Path to input file
        output_file_path: Path to output file
        pattern: Regex pattern for matching lines
        extract_row: Function to extract row data from regex match
    """
    ensure_output_dir(output_file_path)
    
    with open(input_file_path, "r") as input_file, \
         open(output_file_path, 'w', newline='') as output_file:
        writer = csv.writer(output_file)
        
        for line in input_file:
            if match := pattern.search(line):
                row = extract_row(match)
                writer.writerow(row)
