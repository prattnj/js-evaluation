#!/bin/bash

consecutive_newlines=0
expressions=""

while true; do
    read -r -p "" line

    # Check if the line is empty
    if [ -z "$line" ]; then
        ((consecutive_newlines++))
        if [ "$consecutive_newlines" -eq 2 ]; then
            break
        fi
    else
        consecutive_newlines=0
    fi

    # Append the expression to the existing expressions
    expressions="${expressions}${line}\n"
done

# Save the expressions to the source.js file
echo -e "$expressions" > ../source.js
node_modules/.bin/acorn --ecma2024 ../source.js | ruby lib/main.rb
