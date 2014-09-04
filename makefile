install: test
	go install . ./geoelevations
test:
	go test -v ./geoelevations
gofmt:
	gofmt -w . ./geoelevations
goimports:
	goimports -w .
ctags:
	ctags -R .
