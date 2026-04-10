#!/bin/bash
# docs/generate/testkit_examples/testkit_black_box.sh
# Run black-box test example

cd /Users/s3cy/Desktop/go/autoebiten

# Build the state exporter binary first
go build -o ./examples/state_exporter/cmd/state_exporter ./examples/state_exporter/cmd/

# Run the black-box test
go test -v ./examples/testkit/... -run TestPlayerMovement 2>&1
