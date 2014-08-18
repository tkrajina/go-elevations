test:
	go test ./geoelevations
install: test
	go install ./geoelevations
goimports:
	goimports -w .
gofmt:
	gofmt -w .
reload-srtm-data:
	go run reloadjson.go
