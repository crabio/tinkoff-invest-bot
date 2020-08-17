FROM golang:latest

# Add env arguments
ENV GLOBAL_RANK_CSV_FILE_PATH "/data/companies_rank.csv"
ENV MAX_ATTEMPTS 20
ENV CANDLE_INTERVAL "15min"
# DB args
ENV DB_TYPE "postgres"
ENV DB_USER "postgres"
ENV DB_PASSWORD "postgres"
ENV DB_HOSTNAME "timescaledb"
ENV DB_PORT 5432
ENV DB_NAME "tinkoff"

# Create app folder
RUN mkdir /app 

# Add sources
ADD . /app/ 

# Change dir
WORKDIR /app

# Build binary
RUN GOOS=linux CGO_ENABLED=0 go build -o main . 

# Run main app
CMD ["/app/main", "-c", "/config/config.json"]