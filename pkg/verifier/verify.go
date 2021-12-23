//
// Copyright 2021 The Sigstore Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package verifier

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
	"golang.org/x/crypto/openpgp"

	helm_v001 "github.com/sigstore/rekor/pkg/types/helm/v0.0.1"

	"github.com/sigstore/helm-sigstore/pkg/chart"
	"github.com/sigstore/helm-sigstore/pkg/pgp"
)

type Verifier struct {
	ChartManager           *chart.Manager
	Entry                  *helm_v001.V001Entry
	Keyring                openpgp.EntityList
	armoredSignatureReader *io.Reader
	PublicKey              []byte
}

func (v *Verifier) VerifyRekor() error {
	err := v.VerifyChart()

	if err != nil {
		return err
	}

	// Validate Chart Hash
	chartHash, err := v.ChartManager.GetChartDigest()
	rekorChartHash := *v.Entry.HelmObj.Chart.Hash.Value
	if err != nil {
		return errors.Wrap(err, "failed to retrieve Chart Hash")
	}

	if rekorChartHash != chartHash {
		return fmt.Errorf("failed comparing Chart Hash. Value Stored In Rekor Does Not Match Chart. Rekor: '%s'. Chart: '%s", rekorChartHash, chartHash)
	}

	// Verify Public Key Fingerprints and Signature
	rekorFingerprint, err := pgp.GetFingerprintFromPublicKey([]byte(v.Entry.HelmObj.PublicKey.Content))

	if err != nil {
		return errors.New("failed to obtain fingerprint from Rekor public key")
	}

	providedPublicKey, err := pgp.GetFingerprintFromPublicKey(v.PublicKey)
	if err != nil {
		return errors.New("failed to obtain fingerprint from provided public key")
	}

	publicKeyFingerprintCompare := strings.Compare(providedPublicKey, rekorFingerprint)
	if publicKeyFingerprintCompare != 0 {
		return errors.New("failed to comparing Public Key Fingerprints: Value Stored in Rekor Does Not Match Fingerprint of Provided Key")
	}

	armoredSignature, err := ioutil.ReadAll(*v.armoredSignatureReader)
	if err != nil {
		return err
	}

	signatureCompare := bytes.Compare(armoredSignature, []byte(v.Entry.HelmObj.Chart.Provenance.Content))
	if signatureCompare != 0 {
		return errors.New("failed comparing Signature: Value Stored in Rekor Does Not Match Chart Provenance File")
	}

	return nil
}

func (v *Verifier) VerifyChart() error {
	// Get Provenance File
	provenanceFile, err := v.ChartManager.ReadProvenanceFile()

	if err != nil {
		return errors.Wrap(err, "reading Provenance file")
	}

	// Verify Signature by Performing Clearsign
	signer, armoredSignatureReader, err := pgp.VerifySignature(provenanceFile, v.Keyring)

	if err != nil {
		return errors.Wrap(err, "could not verify signature")
	}

	// Set Reader so that it can be used later
	v.armoredSignatureReader = armoredSignatureReader

	publicKey, err := pgp.ExtractPublicKey(signer)

	if err != nil {
		return errors.Wrap(err, "could not extract public key")
	}

	// Set Public key so that it can be used later
	v.PublicKey = publicKey

	return nil
}
