#!/bin/bash
(sleep 45 && cqlsh -f create_tsdb_structure.cql) &
cd / && docker-entrypoint.sh