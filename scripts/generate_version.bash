#!/usr/bin/env bash
set -eu

# Get the version.
version=`git describe --tags HEAD`
# Write out the package.
cat << EOF
package main

var version = "$version"
EOF
