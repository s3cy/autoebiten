#!/bin/bash
# docs/generate/autoui_examples/autoui_waitfor.sh
# Wait for dialog to appear

autoebiten wait-for --condition 'custom:autoui.exists:type=Dialog.found == true' --timeout 5s