FROM golang:1.9.2 as builder
WORKDIR /go/src/github.com/fmitra/dennis
COPY Gopkg.* ./
COPY ./ ./
RUN make deps
RUN make build

FROM scratch
WORKDIR /home
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /go/src/github.com/fmitra/dennis .
EXPOSE 8080
ENTRYPOINT ["./dennis"]
