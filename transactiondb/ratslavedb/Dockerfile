FROM cassandra:latest

RUN mkdir -p /scripts
COPY truncate.cql /scripts/ 
COPY cassandra.yaml /etc/cassandra/
WORKDIR "/"
EXPOSE 9042
CMD [ "/docker-entrypoint.sh" ]