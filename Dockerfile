FROM golang:1.11.1-alpine3.8 as builder

# install git (required by dep ensure)
RUN apk add git

WORKDIR $GOPATH/src/github.com/sinnott74/TodoService
EXPOSE 8000
ENV GO111MODULE on
COPY go.mod go.sum ./
RUN go mod vendor
COPY . ./
RUN CGO_ENABLED=0 go build ./cmd/TodoService/


FROM scratch
WORKDIR /go/
EXPOSE 8000
COPY --from=builder /go/src/github.com/sinnott74/TodoService .
CMD ["./TodoService"]