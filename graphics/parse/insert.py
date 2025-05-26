import csv
import os
import re
from os import path


def parse_insert_to_csv(method, key_kind, load_factor):
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

        with open(input_file_path, "r") as input_file:
            os.makedirs(os.path.dirname(output_file_path), exist_ok=True)

            with open(output_file_path, 'w', newline='') as output_file:
                writer = csv.writer(output_file)

                for line in input_file.readlines():
                    match = pattern.search(line)
                    if match:
                        size = int(match.group('Size'))
                        ns_per_op = float(match.group('NsPerInsert'))
                        b_per_op = int(match.group('BytesPerOp'))

                        writer.writerow([size, ns_per_op, b_per_op])
