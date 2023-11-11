FROM golang:1.21.1 as build

WORKDIR /app

COPY . .

RUN go build -o MNA-Project MNA-project/cmd/api

FROM debian:bookworm

RUN apt-get update -y && apt-get install ca-certificates -y \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=build /app/MNA-Project /app/
COPY --from=build /app/pkg/config/env /app/pkg/config/env

ENV CONFIG_PATH="/app"

RUN useradd -m admin
USER admin

ENTRYPOINT ["./MNA-Project"]
