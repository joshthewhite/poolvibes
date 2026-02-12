FROM golang:1.25 AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /poolvibes .

FROM gcr.io/distroless/static-debian12

COPY --from=build /poolvibes /poolvibes
EXPOSE 8080
ENTRYPOINT ["/poolvibes", "serve", "--addr", ":8080", "--db", "/data/poolvibes.db"]
