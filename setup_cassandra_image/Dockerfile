FROM cassandra:latest

RUN mkdir /scripts
COPY create_tsdb_structure.cql create_database_structure.sh truncate.cql /scripts/
WORKDIR "/scripts"
EXPOSE 9042
CMD ["./create_database_structure.sh"]