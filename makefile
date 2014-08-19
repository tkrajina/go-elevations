test:
	go test -v ./geoelevations
install: test
	go install ./geoelevations
goimports:
	goimports -w .
gofmt:
	gofmt -w ./geoelevations
reload-srtm-data:
	go run reloadjson.go
