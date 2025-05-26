from parse.insert import parse_insert_to_csv
from parse.delete import parse_delete_to_csv
from parse.get import parse_get_to_csv
from parse.other import parse_collisions
from plots.plots import *
from utils.utils import *


def parse_to_csv():
    for method in methods:
        for key_kind in key_kinds:
            for load_factor in load_factors:
                parse_insert_to_csv(method=method, key_kind=key_kind, load_factor=load_factor)
                parse_get_to_csv(method=method, key_kind=key_kind, load_factor=load_factor)
                parse_delete_to_csv(method=method, key_kind=key_kind, load_factor=load_factor)
                parse_collisions(method=method, key_kind=key_kind, load_factor=load_factor)


def main():
    parse_to_csv()

    make_insert_reserve_time_graphics()
    make_insert_no_reserve_time_graphics()
    make_success_get_graphics()
    make_unsuccess_get_graphics()
    make_delete_graphics()
    make_insert_reserve_memory_graphics()


if __name__ == '__main__':
    main()
