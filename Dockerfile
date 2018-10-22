FROM golang:1.11.1-alpine3.8 as builder

RUN apk update && \
  apk add git
# apk add git && \
# apk add ca-certificates && \
# echo "root:x:0:0:root:/root:/bin/sh" > /passwd &&  \
# mkdir /empty

WORKDIR $GOPATH/src/github.com/sinnott74/TodoService
ENV GO111MODULE on

COPY go.mod go.sum ./

RUN go mod download

COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s"

FROM alpine:3.8
EXPOSE 8000

# COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
# COPY --from=builder /passwd /etc/passwd
# COPY --from=builder /bin/sh /bin/sh
# COPY --from=builder /empty /root
COPY --from=builder /go/src/github.com/sinnott74/TodoService/TodoService .
CMD ["./TodoService"]