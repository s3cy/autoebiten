#!/bin/bash
# docs/generate/testkit_examples/testkit_combo_input.sh
# Run combo input test example

cd /Users/s3cy/Desktop/go/autoebiten

# Run the combo input test
go test -v ./examples/testkit/... -run TestComboInput 2>&1
