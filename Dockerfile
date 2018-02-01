FROM scratch

COPY transaction_server /app/
WORKDIR "/app"
EXPOSE 44441
CMD ["./transaction_server"]
