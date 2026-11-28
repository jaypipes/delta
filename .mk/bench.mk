# ==== Benchmark targets ====

##@ Benchmark

.PHONY: bench
bench: bench-cmp ## Run all tests.

.PHONY: test-bench
bench-cmp: ## Run all benchmark tests.
	@cd bench/cmp; go test -bench=. -count=3; cd ../..
