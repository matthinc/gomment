ts-check := tsc --allowJs --checkJs --noEmit --target ES6 --strict

.PHONY: build
build:
	go build

.PHONY: test
test: build
	python3 -m unittest test.api_basic test.admin test.api_comments

.PHONY: tsc
tsc:
	$(ts-check) frontend/gomment.js
	$(ts-check) frontend/admin/gomment-admin.js
