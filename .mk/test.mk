# ==== Test targets ====

##@ Test

.PHONY: test
test: test-unit ## Run all tests.

.PHONY: test-unit
test-unit: ## Run all unit tests.
	@go test -modfile=go_test.mod -v -coverprofile=cover.out ./...

test-bench: ## Run all benchmark tests.
	@cd cmp; go test -modfile=../go_test.mod -bench=. -count=3; cd ..
