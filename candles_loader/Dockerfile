FROM golang:latest

# Add env arguments
ENV SANDBOX_TOKEN ""
ENV PRODUCTION_TOKEN ""
ENV GLOBAL_RANK_CSV_FILE_PATH "/data/companies_rank.csv"
ENV START_LOAD_DATE "2020-01-01T00:00:00.000Z"
ENV MAX_ATTEMPTS 10
# DB args
ENV DB_TYPE "postgres"
ENV DB_USER "postgres"
ENV DB_PASSWORD "postgres"
ENV DB_HOSTNAME "timescaledb"
ENV DB_PORT 5432
ENV DB_NAME "tinkoff"

# Add binary file
ADD build/main /

CMD ["/main"]