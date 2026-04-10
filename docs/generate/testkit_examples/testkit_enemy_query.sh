#!/bin/bash
# docs/generate/testkit_examples/testkit_enemy_query.sh
# Run enemy state query test example

cd /Users/s3cy/Desktop/go/autoebiten

# Build the state exporter binary first
go build -o ./examples/state_exporter/cmd/state_exporter ./examples/state_exporter/cmd/

# Run the enemy query test
go test -v ./examples/testkit/... -run TestEnemyStateQuery 2>&1
