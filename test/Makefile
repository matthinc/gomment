tests: test-image
	python3 system_test_runner.py

.PHONY: test-image
test-image:
	docker build -f test.Dockerfile -t gomment-test ..
