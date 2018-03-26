FROM scratch

COPY transactionServer /app/
WORKDIR "/app"
EXPOSE 44441
CMD ["./transactionServer"]
