.PHONY: build
build:
	go build -o octopiper .

.PHONY: install
install: build
	cp octopiper /usr/local/bin/octopiper
