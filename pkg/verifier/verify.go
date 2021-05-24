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

	"github.com/pkg/errors"
	"github.com/sigstore/helm-sigstore/pkg/chart"
	"github.com/sigstore/helm-sigstore/pkg/pgp"
	helm_v001 "github.com/sigstore/rekor/pkg/types/helm/v0.0.1"
	"golang.org/x/crypto/openpgp"
)

type Verifier struct {
	ChartManager           *chart.ChartManager
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
		return errors.Wrap(err, "Failed to retrieve Chart Hash")
	}

	if rekorChartHash != chartHash {
		return errors.New(fmt.Sprintf("Error Comparing Chart Hash. Value Stored In Rekor Does Not Match Chart. Rekor: '%s'. Chart: '%s", rekorChartHash, chartHash))
	}

	// Verify Public Key and Signature
	publicKeyCompare := bytes.Compare(v.PublicKey, []byte(v.Entry.HelmObj.PublicKey.Content))
	if publicKeyCompare != 0 {
		return errors.New("Error Comparing Public Key: Value Stored in Rekor Does Not Match Key Used to Sign Provenance File")
	}

	armoredSignature, err := ioutil.ReadAll(*v.armoredSignatureReader)
	if err != nil {
		return err
	}

	signatureCompare := bytes.Compare(armoredSignature, []byte(v.Entry.HelmObj.Chart.Provenance.Content))
	if signatureCompare != 0 {
		return errors.New("Error Comparing Signature: Value Stored in Rekor Does Not Match Chart Provenance File")
	}

	return nil
}

func (v *Verifier) VerifyChart() error {

	// Get Provenance File
	provenanceFile, err := v.ChartManager.ReadProvenanceFile()

	if err != nil {
		return errors.Wrap(err, "Failed to Read Provenance file")
	}

	// Verify Signature by Performing Clearsign
	signer, armoredSignatureReader, err := pgp.VerifySignature(provenanceFile, v.Keyring)

	if err != nil {
		return errors.Wrap(err, "Could not verify signature")
	}

	// Set Reader so that it can be used later
	v.armoredSignatureReader = armoredSignatureReader

	publicKey, err := pgp.ExtractPublicKey(signer)

	if err != nil {
		return errors.Wrap(err, "Could not extract public key")
	}

	// Set Public key so that it can be used later
	v.PublicKey = publicKey

	return nil
}
