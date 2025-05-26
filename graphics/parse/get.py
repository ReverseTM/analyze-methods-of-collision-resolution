import csv
import os
import re
from os import path


def parse_get_to_csv(method, key_kind, load_factor):
    for operation in ["UnsuccessGet"]:
        input_file_path = path.join("results", operation, method, "row.txt")
        output_file_path = path.join("data", operation, method, f'{key_kind}_{load_factor}.csv')

        pattern = re.compile(
            rf"Benchmark{re.escape(operation)}/{re.escape(method)}-{re.escape(key_kind)}-{load_factor}-(?P<Size>\d+)-\d+\s+"
            r"\d+\s+"
            r"(?P<NsPerOp>[\d.]+)\s+ns/op\s+"
        )

        with open(input_file_path, "r") as input_file:
            os.makedirs(os.path.dirname(output_file_path), exist_ok=True)

            with open(output_file_path, 'w', newline='') as output_file:
                writer = csv.writer(output_file)

                for line in input_file.readlines():
                    match = pattern.search(line)
                    if match:
                        size = int(match.group('Size'))
                        ns_per_op = float(match.group('NsPerOp'))

                        writer.writerow([size, ns_per_op])
