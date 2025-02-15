#!/bin/bash

for module in storage-operator database-operator message-operator ai-operator common-operator batch-operator; do
    echo "Building $module..."
    cd "$module" || exit
    make docker-build IMG="zeroinfra/$module:latest"
    cd ..
done