#!/bin/bash
cd /
cd - && sleep 40 && cqlsh -f create_tsdb_structure.cql &
docker-entrypoint.sh