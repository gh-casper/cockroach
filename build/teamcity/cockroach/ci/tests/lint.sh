#!/usr/bin/env bash

set -euo pipefail

dir="$(dirname $(dirname $(dirname $(dirname $(dirname "${0}")))))"

source "$dir/teamcity-support.sh"  # For $root
source "$dir/teamcity-bazel-support.sh"  # For run_bazel

tc_prepare

tc_start_block "Run Bazel test"
run_bazel build/teamcity/cockroach/ci/tests/lint_impl.sh
tc_end_block "Run Bazel test"
