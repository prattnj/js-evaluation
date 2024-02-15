#!/bin/bash
while true; do
    read -p "" line

    # Check if the user wants to finish entering expressions
    if [ "$line" == "" ]; then
        break
    fi

    # Append the expression to the existing expressions
    expressions="${expressions}${line}\n"
done

# Save the expressions to the source.js file
echo -e "$expressions" > source.js
node_modules/.bin/acorn --ecma2024 source.js | go run .