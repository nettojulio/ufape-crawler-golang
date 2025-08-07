# Stage 1: Build da aplicação
FROM golang:1.24.4 AS builder

ARG NAME=ufape-crawler-golang

WORKDIR /app

# Copia o código
COPY go.mod ./
COPY go.sum ./
RUN go mod download
RUN go mod tidy

COPY . ./

# Build estático
RUN ./scripts/build.sh

# Stage 2: Distroless
FROM gcr.io/distroless/static-debian11

WORKDIR /

# Copia o binário gerado
ARG NAME=ufape-crawler-golang
COPY --from=builder /app/dist/${NAME}-linux-static /app

# Expõe a porta
EXPOSE 8080

# Comando de entrada
ENTRYPOINT ["/app"]
