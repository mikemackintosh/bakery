.PHONY: pantry
pantry:
	cd pantry && \
		go test -v -tags pantry -run TestGenerate ./... -args -name=$(name)

build:
	go build -o bin/bakery cmd/bakery.go

build-complete:
	rice embed-go
	make build
