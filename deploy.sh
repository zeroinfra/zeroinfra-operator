#!/bin/bash

for module in storage-operator database-operator message-operator ai-operator common-operator batch-operator; do
    echo "Deploying $module..."
    cd "$module" || exit
    make deploy IMG="zeroinfra/$module:latest"
    cd ..
done