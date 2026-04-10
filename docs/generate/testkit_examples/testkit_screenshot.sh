#!/bin/bash
# docs/generate/testkit_examples/testkit_screenshot.sh
# Run screenshot capture test example

cd /Users/s3cy/Desktop/go/autoebiten

# Build the state exporter binary first
go build -o ./examples/state_exporter/cmd/state_exporter ./examples/state_exporter/cmd/

# Run the screenshot test
go test -v ./examples/testkit/... -run TestScreenshotCapture 2>&1
