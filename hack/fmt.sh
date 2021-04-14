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

export CGO_ENABLED=0
export GO111MODULE=on
export GOFLAGS="-mod=vendor"

TARGETS="$@"

if [ -n "$TARGETS" ]; then
    echo "Running reimport.py"
    cmd="reimport3.py ${REPO_PKG} ${TARGETS}"
    $cmd
    echo

    echo "Running goimports:"
    cmd="goimports -w ${TARGETS}"
    echo "$cmd"
    $cmd
    echo

    echo "Running gofmt:"
    cmd="gofmt -s -w ${TARGETS}"
    echo "$cmd"
    $cmd
    echo
fi

echo "Running shfmt:"
cmd="find . -path ./vendor -prune -o -name '*.sh' -exec shfmt -l -w -ci -i 4 {} \;"
echo "$cmd"
eval "$cmd" # xref: https://stackoverflow.com/a/5615748/244009
echo
