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

package cmd

import (
	"bytes"
	"fmt"

	"encoding/base64"

	"github.com/go-openapi/runtime"
	"github.com/pkg/errors"
	"github.com/sigstore/helm-sigstore/pkg/chart"
	"github.com/sigstore/helm-sigstore/pkg/pgp"
	"github.com/sigstore/helm-sigstore/pkg/rekor"
	"github.com/sigstore/helm-sigstore/pkg/types"
	"github.com/sigstore/helm-sigstore/pkg/verifier"
	"github.com/sigstore/rekor/pkg/generated/models"
	rekortypes "github.com/sigstore/rekor/pkg/types"
	helm_v001 "github.com/sigstore/rekor/pkg/types/helm/v0.0.1"
	"github.com/spf13/cobra"
)

func NewVerifyCmd() *cobra.Command {

	verifyOptions := types.CLIOptions{}

	// searchCmd represents the upload command
	verifyCmd := &cobra.Command{
		Use:   "verify [PATH_TO_PACKAGED_CHART]",
		Short: "Verify a Signed Helm Chart",
		RunE: func(cmd *cobra.Command, args []string) error {

			if len(args) != 1 {
				return errors.New("1 argument (Path to packaged chart) is required")
			}

			chartPath := args[0]
			publicKeyPath := verifyOptions.PublicKey
			keyringPath := verifyOptions.Keyring
			rekorServer := verifyOptions.RekorServer

			chartManager, err := chart.NewChartManager(chartPath)

			if err != nil {
				return err
			}

			rekor, err := rekor.NewRekor(rekorServer)

			if err != nil {
				return err
			}

			uuids, err := rekor.Search(chartManager)

			if err != nil {
				return err
			}

			if len(uuids) == 0 {
				return errors.New("Unable to verify Chart: No Record Found")
			}

			uuid := uuids[len(uuids)-1]

			logEntryAnon, err := rekor.GetByUUID(uuid)

			if err != nil {
				return errors.Wrapf(err, "Error Retriving Log Entry: %s", uuid)
			}

			b, err := base64.StdEncoding.DecodeString(logEntryAnon.Body.(string))
			if err != nil {
				return err
			}

			pe, err := models.UnmarshalProposedEntry(bytes.NewReader(b), runtime.JSONConsumer())
			if err != nil {
				return err
			}
			eimpl, err := rekortypes.NewEntry(pe)
			if err != nil {
				return err
			}

			helmEntry, ok := eimpl.(*helm_v001.V001Entry)
			if !ok {
				return errors.New("cannot unmarshal non Helm v0.0.1 type")
			}

			// Get Keyring
			keyring, err := pgp.GetKeyring(keyringPath, publicKeyPath)

			if err != nil {
				return errors.Wrap(err, "Could not retrieve keyring")
			}

			verifier := verifier.Verifier{
				ChartManager: chartManager,
				Keyring:      keyring,
				Entry:        helmEntry,
			}

			err = verifier.VerifyRekor()

			if err != nil {
				return err
			}

			fmt.Println("Chart Verified Successfully From Helm entry:")
			fmt.Printf("\nRekor Server: %s", rekorServer)
			fmt.Printf("\nRekor Index: %d", int(*logEntryAnon.LogIndex))
			fmt.Printf("\nRekor UUID: %s\n", uuid)

			return nil
		},
	}

	addRekorFlags(verifyCmd, &verifyOptions)
	addPkiFlags(verifyCmd, &verifyOptions)

	return verifyCmd
}
