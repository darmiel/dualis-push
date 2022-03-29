FROM golang:1.18

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod graph | awk '{if ($1 !~ "@") print $2}' | xargs go get
COPY . .
RUN go build -o /check
ENTRYPOINT [ "/check" ]