FROM golang

WORKDIR /app

ARG target

COPY go.mod go.sum ./
RUN go mod download

COPY ./api ./api
COPY ./books ./books
COPY ./cmd/${target} ./cmd/${target}

RUN go build -o app ./cmd/${target}

EXPOSE 8080

ENTRYPOINT ["./app"]
