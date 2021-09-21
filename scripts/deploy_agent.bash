#!/usr/bin/env bash
set -eu
flyctl deploy --config dd.fly.toml --dockerfile ./Datadog.dockerfile
