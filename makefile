test:
	go test -v ./geoelevations

gofmt:
	gofmt -w . ./geoelevations

generate-urls:
	go run cmd/generate_urls/generate_urls.go
	gofmt -w . ./geoelevations

goimports:
	goimports -w .

ctags:
	ctags -R .
