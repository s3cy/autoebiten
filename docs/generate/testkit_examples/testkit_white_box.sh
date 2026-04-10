#!/bin/bash
# docs/generate/testkit_examples/testkit_white_box.sh
# Run white-box test example

cd /Users/s3cy/Desktop/go/autoebiten

go test -v ./examples/testkit/... -run TestPlayerMovesRight 2>&1
