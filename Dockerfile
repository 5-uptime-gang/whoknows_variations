FROM golang:1.25

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# build the main app from /cmd
RUN go build -o /usr/local/bin/whoknows_variations ./cmd

CMD ["whoknows_variations"]