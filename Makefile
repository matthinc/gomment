build:
	go build

.PHONY: test
test: build
	python3 ./test/main.py
