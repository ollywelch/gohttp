FROM golang:1.20 as builder

ENV CGO_ENABLED=0 \
    GOOS=linux

WORKDIR /root

COPY go.mod go.sum ./

RUN go mod download

COPY *.go ./

RUN go build -o /gohttp

FROM scratch

COPY --from=builder /gohttp /gohttp

EXPOSE 3000

CMD [ "/gohttp" ]
