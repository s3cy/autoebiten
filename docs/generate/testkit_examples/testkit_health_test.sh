#!/bin/bash
# docs/generate/testkit_examples/testkit_health_test.sh
# Run health modification test example

cd /Users/s3cy/Desktop/go/autoebiten

# Build the state exporter binary first
go build -o ./examples/state_exporter/cmd/state_exporter ./examples/state_exporter/cmd/

# Run the health test
go test -v ./examples/testkit/... -run TestHealthModification 2>&1
