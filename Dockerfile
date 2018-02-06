FROM golang:1.9.2 as builder
WORKDIR /go/src/github.com/fmitra/dennis
COPY Gopkg.* ./
RUN go get -v -u github.com/golang/dep/cmd/dep && dep ensure -vendor-only -v
COPY ./ ./
RUN make build

FROM scratch
WORKDIR /home
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/github.com/fmitra/dennis .
EXPOSE 8080
ENTRYPOINT ["./dennis", "--config=/etc/dennis/config.json"]
