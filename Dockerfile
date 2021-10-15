FROM golang:1.16-buster as build

WORKDIR /build

COPY go.mod go.sum main.go ./
COPY utils* ./utils/

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o main_go .

FROM scratch AS production

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /

COPY --from=build /build/main_go .
COPY app.env .

CMD ["./main_go"]
