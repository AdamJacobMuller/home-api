VERSION=$(shell git describe  --tags)
BUILD_TIME=$(shell date -u +%FT%H:%M:%SZ)

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

run: templates
	go run ${LDFLAGS} cli/api/main.go -configuration=config.json

templates: api/server/templates/*.qtpl
	qtc -dir api/server/templates
