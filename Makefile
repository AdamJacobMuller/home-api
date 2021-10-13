VERSION=$(shell git describe  --tags)
BUILD_TIME=$(shell date -u +%FT%H:%M:%SZ)

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}"

linux:
	GOOS=linux go build -o build/api-linux ${LDFLAGS} cli/api/main.go
	scp build/api-linux adam@nighe.adam.gs:~/home-api
	ssh adam@10.0.8.3 /usr/bin/sudo /bin/systemctl restart home-api

run:
	go run ${LDFLAGS} cli/api/main.go -configuration=config.json

run-%:
	go run ${LDFLAGS} cli/api/main.go -configuration=$*.json
