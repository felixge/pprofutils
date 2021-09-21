#!/usr/bin/env bash
set -eu

# Get the version.
version=`git describe --tags HEAD`
# Write out the package.
cat << EOF
// Code generated ./scripts/generate_version.bash DO NOT EDIT.

package main

var version = "$version"
EOF
