#!/bin/bash
# docs/generate/testkit_examples/testkit_damage_test.sh
# Run white-box damage test example

cd /Users/s3cy/Desktop/go/autoebiten

# Run the white-box damage test
go test -v ./examples/testkit/... -run TestPlayerTakesDamage 2>&1
