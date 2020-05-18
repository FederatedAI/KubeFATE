FROM golang as build

WORKDIR /data/projects/fatecloud

COPY . ./

RUN go build -a -ldflags '-s' -installsuffix cgo -o kubefate kubefate.go

FROM centos

WORKDIR /data/projects/fatecloud

COPY --from=0 /data/projects/fatecloud/kubefate /data/projects/fatecloud/
COPY --from=0 /data/projects/fatecloud/config.yaml /data/projects/fatecloud/

EXPOSE 8080

CMD ["./kubefate service"]

ENTRYPOINT ["/bin/sh","-c"]

