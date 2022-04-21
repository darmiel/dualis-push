FROM golang:1.18 AS builder

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod graph | awk '{if ($1 !~ "@") print $2}' | xargs go get
COPY . .
RUN CGO_ENABLED=0 go build -o /check

FROM alpine:latest
WORKDIR /usr/src/app

COPY --from=builder /check .
# copy config files
COPY *.toml .

RUN ls -larth

ENTRYPOINT [ "./check" ]