#!/bin/bash
cd /
cd - && sleep 45 && cqlsh -f create_tsdb_structure.cql &
docker-entrypoint.sh