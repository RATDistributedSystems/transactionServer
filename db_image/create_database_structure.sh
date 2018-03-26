#!/bin/bash
(sleep 65 && cqlsh -f create_tsdb_structure.cql) &
cd / && docker-entrypoint.sh