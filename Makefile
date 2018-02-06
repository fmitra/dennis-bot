build:
	CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -ldflags '-s' -o dennis .
