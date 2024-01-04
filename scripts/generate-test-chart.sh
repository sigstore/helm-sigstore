#!/usr/bin/env bash
# Copyright 2024 The Sigstore Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

echo "Creating throw away key..."
mkdir -p .gnupg
gpg --batch --passphrase '' --quick-generate-key "helm-sigstore-test"
gpg --export-secret-keys > .gnupg/sigstore-secring.gpg

echo "Creating, packaging and signing chart temp chart..."
helm create helm-sigstore-test
helm package --sign --key 'helm-sigstore-test' --keyring .gnupg/sigstore-secring.gpg helm-sigstore-test
cat helm-sigstore-test-0.1.0.tgz.prov
