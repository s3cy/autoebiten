#!/bin/bash
# docs/generate/autoui_examples/autoui_xpath_level.sh
# XPath query for widgets with numeric level > 40

autoebiten custom autoui.xpath --request "//*[number(@level) > 40]"