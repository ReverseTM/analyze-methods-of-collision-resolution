.PHONY: bench all

RESULTS_DIR := results

# Ops: InsertReserve InsertNoReserve, SuccessGet, UnsuccessGet, Delete

# Параметры: OP=InsertNoReserve, METHOD=RobinHood и т.п.
bench:
	@mkdir -p $(RESULTS_DIR)
	go test -bench="Benchmark$(OP)$$" -run=^$$(echo) ./... \
		| tee $(RESULTS_DIR)/$(OP).txt
