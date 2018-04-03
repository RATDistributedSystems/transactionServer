#!/bin/bash
mv cassandra.yaml /etc/cassandra/
(sleep 65 && cqlsh -f create_tsdb_structure.cql) &
cd / && docker-entrypoint.sh