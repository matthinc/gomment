build:
	go build

.PHONY: test
test: build
	python -m unittest test.api_basic test.admin test.api_comments
