#!/bin/bash

# Exit if any command fails
set -e

# Fetch all tags and sort them to get the latest
git fetch --tags

# Get the latest tag (considering semver sorting)
latest_tag=$(git tag -l | sort -V | tail -n1)

# If there are no tags yet, start with 0.0.0
if [ -z "$latest_tag" ]; then
    latest_tag="0.0.0"
fi

# Split the tag into major, minor, and patch numbers
IFS='.' read -ra parts <<< "$latest_tag"
major=${parts[0]}
minor=${parts[1]}
patch=${parts[2]}

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

echo "tag=$new_tag" >> $GITHUB_OUTPUT
