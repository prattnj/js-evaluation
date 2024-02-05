#!/bin/bash
read -p $'Enter javascript arithmetic expression:\n' expression
json_object=`echo "$expression" | node_modules/.bin/acorn --ecma2024`
go run . "$json_object"