FROM scratch

COPY transactionserver /app/
WORKDIR "/app"
EXPOSE 44441
CMD ["./transactionserver"]
