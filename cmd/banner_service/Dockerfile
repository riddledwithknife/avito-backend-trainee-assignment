FROM golang:alpine AS builder-banner-service

WORKDIR /avito-backend-trainee-assignment

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/banner_service/banner_service.go ./cmd/banner_service/
COPY configs/config.json ./configs/
COPY internal/config/config.go ./internal/config/
COPY internal/handlers/handlers.go ./internal/handlers/
COPY internal/models/models.go ./internal/models/

RUN go build -o ./banner_service ./cmd/banner_service/banner_service.go

FROM alpine AS runner-banner-service

COPY --from=builder-banner-service /avito-backend-trainee-assignment/banner_service /
COPY --from=builder-banner-service /avito-backend-trainee-assignment/configs/config.json /

CMD ["/banner_service"]