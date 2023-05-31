#!/bin/bash

function main() {
     if [ -n "$1" ] && [ "${1}" = "github" ]; then
         echo "VERSION=$(printVersion) >> $GITHUB_OUTPUT"
     else
         printVersion
     fi
}

function printVersion() {
    git describe --always
}

main "$@"
