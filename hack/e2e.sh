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

export CGO_ENABLED=1
export GO111MODULE=on
export GOFLAGS="-mod=vendor"

GINKGO_ARGS=${GINKGO_ARGS:-}
TEST_ARGS=${TEST_ARGS:-}
DOCKER_REGISTRY=${DOCKER_REGISTRY:-}

echo "Running e2e tests:"
cmd="ginkgo -r --v -race --progress --trace --noisyPendings=false ${GINKGO_ARGS} test -- --docker-registry=${DOCKER_REGISTRY} ${TEST_ARGS}"
echo $cmd
$cmd
