FROM golang:1.17.6-buster AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY pkg pkg
COPY cmd/greeter cmd/greeter

RUN go build -o greeter ./cmd/greeter/

FROM gcr.io/distroless/base-debian10

WORKDIR /
COPY --from=build /app/greeter /bin/

EXPOSE 8080
USER nonroot:nonroot

ENTRYPOINT ["greeter"]
