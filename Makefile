VERSION=$(shell git describe  --tags)
BUILD_TIME=$(shell date -u +%FT%H:%M:%SZ)

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

linux: templates
	GOOS=linux go build -o build/api-linux ${LDFLAGS} cli/api/main.go

run: templates
	go run ${LDFLAGS} cli/api/main.go -configuration=config.json

run-%: templates
	go run ${LDFLAGS} cli/api/main.go -configuration=$*.json

templates: api/server/templates/*.qtpl
	qtc -dir api/server/templates
