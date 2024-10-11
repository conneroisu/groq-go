# -include .env
#  file: Makefile
#  url: https://github.com/conneroisu/groq-go/Makefile
#  description: Makefile for the project

export MAKEFLAGS += --always-make --print-directory
SHELLFLAGS = -e
.PHONY: tidy
tidy:
	@sh ./scripts/makefile/tidy.sh

.PHONY: dev
dev:
	@sh ./scripts/makefile/dev.sh
	@make clean

.PHONY: vet
vet:
	@sh ./scripts/makefile/vet.sh

.PHONY: lint
lint:
	@+sh ./scripts/makefile/lint.sh

.PHONY: test
test:
	@+sh ./scripts/makefile/test.sh

.PHONY: install
install:
	@sh ./scripts/makefile/install.sh

.PHONY: clean
clean:
	@sh ./scripts/makefile/clean.sh

.PHONY: coverage
coverage:
	@sh ./scripts/makefile/coverage.sh

.PHONY: fmt
fmt:
	@sh ./scripts/makefile/fmt.sh

.PHONY: docs
docs:
	@sh ./scripts/makefile/docs.sh
