develop:
	cp config/config.example.json config/config.json
	cp docker-compose.example.yml docker-compose.yml

dev_dependencies:
	go get -v -u honnef.co/go/tools/cmd/megacheck
	go get -v -u golang.org/x/lint/golint
	go get -v -u github.com/golang/dep/cmd/dep && dep ensure -vendor-only -v

dependencies:
	go get -v -u github.com/golang/dep/cmd/dep && dep ensure -vendor-only -v

test_and_lint:
	go fmt ./...
	go vet ./...
	megacheck $$(go list ./...)
	golint $$(go list ./...)
	go test -v -cover ./...

build:
	CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -ldflags '-s' -o dennis-bot ./cmd/dennis-bot
