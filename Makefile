ts-check := tsc --allowJs --checkJs --noEmit --target ES6 --strict

.PHONY: build
build:
	go build

.PHONY: test-unit
test-unit:
	go test -v ./...

.PHONY: test-system
test-system: build
	python3 -m unittest test.api_basic test.admin test.api_comments

.PHONY: tsc
tsc:
	$(ts-check) frontend/gomment.js
	$(ts-check) frontend/admin/gomment-admin.js

.PHONY: test
test: test-unit tsc test-system
