FROM golang:1.20 as builder
COPY . /app
WORKDIR /app
ARG LOG_DIR=/app/logs
RUN mkdir -p ${LOG_DIR}
# Environment Variables
ENV LOG_FILE_LOCATION=${LOG_DIR}/app.log 
COPY go.mod go.sum ./
RUN go mod download
COPY . /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /bin/app ./cmd/app

FROM scratch
COPY --from=builder /app/config /config
COPY --from=builder /bin/app /app
# Declare volumes to mount
#VOLUME [${LOG_DIR}]
EXPOSE 8080
CMD ["/app"]