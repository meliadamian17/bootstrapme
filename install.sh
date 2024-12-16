#!/usr/bin/env bash

set -e

go build -o bootstrapme ./cmd/bootstrapme

sudo mv bootstrapme /usr/local/bin/bootstrapme

mkdir -p ~/.config/bootstrapme

cp -r defaults/* ~/.config/bootstrapme/

echo "bootstrapme installed successfully!"
echo "Default configs copied to ~/.config/bootstrapme/"
echo "Run 'bootstrapme' to start."

