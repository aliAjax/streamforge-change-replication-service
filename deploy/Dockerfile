FROM golang:1.23 AS build
WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -o /out/streamforge-cdc ./cmd/streamforge-cdc
FROM gcr.io/distroless/static-debian12
COPY --from=build /out/streamforge-cdc /streamforge-cdc
EXPOSE 8092
ENTRYPOINT ["/streamforge-cdc"]
