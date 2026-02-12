FROM golang:1.25 AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /poolio .

FROM gcr.io/distroless/static-debian12

COPY --from=build /poolio /poolio
EXPOSE 8080
ENTRYPOINT ["/poolio", "serve", "--addr", ":8080", "--db", "/data/poolio.db"]
