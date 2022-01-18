FROM golang:1.17 as builder

WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY kubefate.go kubefate.go
COPY pkg/ pkg/
COPY docs/docs.go docs/docs.go 
COPY config.yaml config.yaml

ARG LDFLAGS
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -ldflags '-s' -installsuffix cgo -o kubefate kubefate.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/kubefate .
COPY --from=builder /workspace/config.yaml  .

USER nonroot:nonroot

EXPOSE 8080

CMD ["service"]

ENTRYPOINT ["/kubefate"]
