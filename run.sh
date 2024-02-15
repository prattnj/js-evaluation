#!/bin/bash
read -p $'Enter javascript arithmetic expression:\n' expression
"$expression" | node_modules/.bin/acorn --ecma2024 | go run .