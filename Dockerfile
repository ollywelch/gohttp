FROM golang:1.20 as builder

ENV CGO_ENABLED=0

WORKDIR /root

COPY . .

RUN go build -o /gohttp

FROM scratch

COPY --from=builder /gohttp /gohttp

EXPOSE 3000

CMD [ "/gohttp" ]
