#!/bin/bash
# docs/generate/autoui_examples/autoui_settext.sh
# Set text in a TextInput widget

autoebiten custom autoui.call --request '{"target":"id=name-input","method":"SetText","args":["Alice"]}'