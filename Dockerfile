FROM scratch

COPY transactionServer frontend /app/
WORKDIR "/app"
EXPOSE 44440
CMD ["./transactionServer"]
