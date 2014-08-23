install: gofmt
	go install . ./geoelevations
gofmt:
	gofmt -w .
goimports:
	goimports -w .
test: install
	go test -v ./geoelevations
ctags:
	ctags -R .
