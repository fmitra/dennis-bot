deps:
	go get -v -u github.com/golang/dep/cmd/dep && dep ensure -vendor-only -v

build:
	CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -ldflags '-s' -o dennis .
