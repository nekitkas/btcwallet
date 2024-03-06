# build Stage
FROM golang:alpine AS build

# installing build dependencies
RUN apk add --no-cache gcc musl-dev make

WORKDIR /usr/local/src

COPY ./ ./
RUN CGO_ENABLED=1 go build -o ./build/api cmd/api/main.go

# final Stage
FROM alpine:latest

COPY --from=build /usr/local/src/build/ /
COPY --from=build /usr/local/src/configs/config.json configs/config.json
COPY --from=build /usr/local/src/internal/db/database.db internal/db/database.db
COPY --from=build /usr/local/src/migrations/ migrations/

EXPOSE 8080
CMD ["/api"]