FROM cassandra:latest

RUN mkdir -p /scripts
COPY truncate.cql start.sh cassandra.yaml /scripts/ 
WORKDIR "/"
EXPOSE 9042
CMD [ "start.sh" ]