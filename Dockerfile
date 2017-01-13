FROM scratch

ADD ca-certificates.crt /etc/ssl/certs/

#CGO_ENABLED=0 GOOS=linux go build -a --ldflags="-s" --installsuffix cgo -o go_million

ADD go_million /

ENTRYPOINT ["/go_million"]

