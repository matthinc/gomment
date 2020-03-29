build:
	go build

.PHONY: test
test: build
	python3 -m unittest test.api_basic test.admin test.api_comments
