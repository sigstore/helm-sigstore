#!/usr/bin/env bats

#
# Copyright 2024 The Sigstore Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

@test "print help" {
  cmd="bin/helm-sigstore help"
  run ${cmd}

  echo "${cmd} : ${status} : ${output}"
  [ "$status" -eq 0 ]
}

@test "print version" {
  cmd="bin/helm-sigstore version"
  run ${cmd}

  echo "${cmd} : ${status} : ${output}"
  [ "$status" -eq 0 ]
}

@test "upload packaged chart" {
  cmd="bin/helm-sigstore upload helm-sigstore-test-0.1.0.tgz --keyring .gnupg/sigstore-secring.gpg"
  run ${cmd}

  echo "${cmd} : ${status} : ${output}"
  [ "$status" -eq 0 ]
}

@test "search packaged chart" {
  cmd="bin/helm-sigstore search helm-sigstore-test-0.1.0.tgz"
  run ${cmd}

  echo "${cmd} : ${status} : ${output}"
  [ "$status" -eq 0 ]
}

@test "verify packaged chart" {
  cmd="bin/helm-sigstore verify helm-sigstore-test-0.1.0.tgz --keyring .gnupg/sigstore-secring.gpg"
  run ${cmd}

  echo "${cmd} : ${status} : ${output}"
  [ "$status" -eq 0 ]
}
