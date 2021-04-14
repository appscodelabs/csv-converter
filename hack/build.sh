#!/usr/bin/env bash

# Copyright AppsCode Inc. and Contributors
#
# This program is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# This program is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with this program.  If not, see <https://www.gnu.org/licenses/agpl-3.0.txt>.

set -eou pipefail

if [ -z "${OS:-}" ]; then
    echo "OS must be set"
    exit 1
fi
if [ -z "${ARCH:-}" ]; then
    echo "ARCH must be set"
    exit 1
fi
if [ -z "${VERSION:-}" ]; then
    echo "VERSION must be set"
    exit 1
fi

export CGO_ENABLED=0
export GOARCH="${ARCH}"
export GOOS="${OS}"
export GO111MODULE=on
export GOFLAGS="-mod=vendor"

go install \
    -installsuffix "static" \
    -ldflags "                                          \
      -X main.Version=${VERSION}                        \
      -X main.VersionStrategy=${version_strategy:-}     \
      -X main.GitTag=${git_tag:-}                       \
      -X main.GitBranch=${git_branch:-}                 \
      -X main.CommitHash=${commit_hash:-}               \
      -X main.CommitTimestamp=${commit_timestamp:-}     \
      -X main.GoVersion=$(go version | cut -d " " -f 3) \
      -X main.Compiler=$(go env CC)                     \
      -X main.Platform=${OS}/${ARCH}                    \
    " \
    ./...
