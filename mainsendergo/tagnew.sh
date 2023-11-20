#!/bin/bash

# Exit if any command fails
set -e

# Fetch all tags
git fetch --tags

# Get the latest tag
latest_tag=$(git describe --tags `git rev-list --tags --max-count=1`)

# If there are no tags yet, start with 0.1.0
if [ -z "$latest_tag" ]; then
    latest_tag="0.1.0"
fi

# Split the tag into major, minor, and patch numbers
IFS='.' read -ra ADDR <<< "$latest_tag"
major=${ADDR[0]}
minor=${ADDR[1]}
patch=${ADDR[2]}

# Increment the patch number
new_patch=$((patch + 1))

# Create new tag
new_tag="$major.$minor.$new_patch"

echo "Creating and pushing new tag: $new_tag"

# Create the new tag
git tag $new_tag

# Push the new tag
git push origin $new_tag

echo "Tag $new_tag created and pushed successfully."
