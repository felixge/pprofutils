#!/usr/bin/env bash
set -eu
go generate ./cmd/pprofutils/
flyctl deploy
