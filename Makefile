.PHONY: pantry
pantry:
	cd pantry && \
		go test -v -tags pantry -run TestGenerate ./... -args -name=$(name)
