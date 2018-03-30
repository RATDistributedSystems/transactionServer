FROM scratch

COPY transactionServer frontend /app/
WORKDIR "/app"
EXPOSE 44441
CMD ["./transactionServer"]
