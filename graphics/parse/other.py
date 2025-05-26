import csv
import os
from os import path


def parse_collisions(method, key_kind, load_factor):
    metric = "Collisions"

    input_file_path = path.join("results", metric, method, load_factor, f'{key_kind}.csv')
    output_file_path = path.join("data", metric, method, f'{key_kind}_{load_factor}.csv')

    with open(input_file_path, "r", newline='') as input_file:
        os.makedirs(os.path.dirname(output_file_path), exist_ok=True)

        with open(output_file_path, 'w', newline='') as output_file:
            reader = csv.reader(input_file)
            writer = csv.writer(output_file)

            for line in reader:
                writer.writerow(line)
