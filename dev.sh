#!/usr/bin/env sh

# Compile application and run with provided arguments
set -e

PROGRAM_NAME="bugout"

go build -o "$PROGRAM_NAME" ./cmd/bugout

./"$PROGRAM_NAME" "$@"
