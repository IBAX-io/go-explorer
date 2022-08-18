#!/bin/bash
set -e -x

HOMEDIR=$(pwd)

function buildpkg() {
    buildBin=$1
    buildModule=$2
    buildFile=$3
    buildBranch=$(git rev-parse --abbrev-ref HEAD 2>/dev/null || echo unknown)
    buildDate=$(date -u "+%Y-%m-%d-%H:%M:%S(UTC)")
    commitHash=$(git rev-parse --short HEAD 2>/dev/null || echo unknown)
    go build -o "$buildBin" -ldflags "-X $buildModule/models.buildBranch=$buildBranch -X $buildModule/models.buildDate=$buildDate -X $buildModule/models.commitHash=$commitHash" "$buildFile"
}

buildpkg go-explorer "github.com/IBAX-io/go-explorer" "$HOMEDIR/main.go"
