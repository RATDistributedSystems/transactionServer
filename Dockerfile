FROM scratch

COPY transactionserver config.json /app/
WORKDIR "/app"
EXPOSE 44441
CMD ["./transactionserver"]
