from utils.utils import *


def make_insert_no_reserve_time_graphics():
    input_dir = path.join("data", "InsertNoReserve")
    output_dir = path.join("graphics", "InsertNoReserve")

    for key_kind in key_kinds:
        for load_factor in load_factors:
            if load_factor == "0.40":
                make_graphic(
                    input_dir,
                    output_dir,
                    f'{key_kind}_{load_factor}',
                    "Количество элементов",
                    "Среднее время 1 операции (ns)",
                    [0, 1]
                )
            else:
                make_two_graphic(
                    input_dir,
                    output_dir,
                    f'{key_kind}_{load_factor}',
                    "Количество элементов",
                    "Среднее время 1 операции (ns)",
                    [0, 1]
                )


def make_insert_reserve_time_graphics():
    input_dir = path.join("data", "InsertReserve")
    output_dir = path.join("graphics", "InsertReserve")

    for key_kind in key_kinds:
        for load_factor in load_factors:
            make_graphic(
                input_dir,
                output_dir,
                f'{key_kind}_{load_factor}',
                "Количество элементов",
                "Среднее время 1 операции (ns)",
                [0, 1]
            )


def make_insert_reserve_memory_graphics():
    input_dir = path.join("data", "InsertNoReserve")
    output_dir = path.join("graphics", "AllocateMemoryInsertNoReserve")

    for key_kind in key_kinds:
        for load_factor in load_factors:
            if load_factor == "0.40":
                make_graphic(
                    input_dir,
                    output_dir,
                    f'{key_kind}_{load_factor}',
                    "Количество элементов",
                    "Количество выделенной памяти (bytes)",
                    [0, 2]
                )
            else:
                make_two_graphic(
                    input_dir,
                    output_dir,
                    f'{key_kind}_{load_factor}',
                    "Количество элементов",
                    "Количество выделенной памяти (bytes)",
                    [0, 2]
                )


def make_success_get_graphics():
    input_dir = path.join("data", "SuccessGet")
    output_dir = path.join("graphics", "SuccessGet")

    for key_kind in key_kinds:
        for load_factor in load_factors:
            make_graphic(
                input_dir,
                output_dir,
                f'{key_kind}_{load_factor}',
                "Количество элементов",
                "Среднее время 1 операции (ns)",
                [0, 1]
            )


def make_unsuccess_get_graphics():
    input_dir = path.join("data", "UnsuccessGet")
    output_dir = path.join("graphics", "UnsuccessGet")

    for key_kind in key_kinds:
        for load_factor in load_factors:
            make_graphic(
                input_dir,
                output_dir,
                f'{key_kind}_{load_factor}',
                "Количество элементов",
                "Среднее время 1 операции (ns)",
                [0, 1]
            )


def make_delete_graphics():
    input_dir = path.join("data", "Delete")
    output_dir = path.join("graphics", "Delete")

    for key_kind in key_kinds:
        for load_factor in load_factors:
            make_graphic(
                input_dir,
                output_dir,
                f'{key_kind}_{load_factor}',
                "Количество элементов",
                "Среднее время 1 операции (ns)",
                [0, 1]
            )
