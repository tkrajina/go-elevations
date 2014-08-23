install: gofmt
	go install . ./geoelevations
gofmt:
	gofmt -w . ./geoelevations
goimports:
	goimports -w .
test: install
	go test -v ./geoelevations
ctags:
	ctags -R .
