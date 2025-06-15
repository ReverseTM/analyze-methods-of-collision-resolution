from plots.plots import *
from parse.parser import *


def main():
    parse_results()

    make_insert_reserve_time_graphics()
    make_insert_no_reserve_time_graphics()
    make_success_get_graphics()
    make_unsuccess_get_graphics()
    make_delete_graphics()
    make_insert_reserve_memory_graphics()


if __name__ == '__main__':
    main()
