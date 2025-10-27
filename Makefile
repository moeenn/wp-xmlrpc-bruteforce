run:
	go run src/*.go


build:
	go build -o exploit src/*.go


.PHONY: run
