VERSION=$(shell git describe  --tags)
BUILD_TIME=$(shell date -u +%FT%H:%M:%SZ)

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

linux: templates
	GOOS=linux go build -o build/api-linux ${LDFLAGS} cli/api/main.go
	scp build/api-linux adam@nighe.adam.gs:~/api-linux.1
	ssh adam@nighe.adam.gs supervisorctl -s http://127.0.0.1:8081 stop home-api
	ssh adam@nighe.adam.gs mv api-linux api-linux.$(shell date +%s)
	ssh adam@nighe.adam.gs mv api-linux.1 api-linux
	ssh adam@nighe.adam.gs supervisorctl -s http://127.0.0.1:8081 start home-api

run: templates
	go run ${LDFLAGS} cli/api/main.go -configuration=config.json

run-%: templates
	go run ${LDFLAGS} cli/api/main.go -configuration=$*.json

templates: api/server/templates/*.qtpl
	qtc -dir api/server/templates
