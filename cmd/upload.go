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
	"fmt"

	"github.com/go-openapi/swag"
	"github.com/pkg/errors"
	"github.com/sigstore/helm-sigstore/pkg/chart"
	"github.com/sigstore/helm-sigstore/pkg/pgp"
	"github.com/sigstore/helm-sigstore/pkg/rekor"
	"github.com/sigstore/helm-sigstore/pkg/types"
	"github.com/sigstore/helm-sigstore/pkg/verifier"
	"github.com/spf13/cobra"
)

func NewUploadCmd() *cobra.Command {

	uploadOptions := types.CLIOptions{}

	// uploadCmd represents the upload command
	uploadCmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload Signed Helm Chart",
		RunE: func(cmd *cobra.Command, args []string) error {

			if len(args) != 1 {
				return errors.New("1 argument (Path to packaged chart) is required")
			}

			chartPath := args[0]

			chartManager, err := chart.NewChartManager(chartPath)

			if err != nil {
				return err
			}

			r, err := rekor.NewRekor(uploadOptions.RekorServer)

			if err != nil {
				return err
			}

			publicKeyPath := uploadOptions.PublicKey
			keyringPath := uploadOptions.Keyring
			rekorServer := uploadOptions.RekorServer

			// Get Keyring
			keyring, err := pgp.GetKeyring(keyringPath, publicKeyPath)

			if err != nil {
				return errors.Wrap(err, "Could not retrieve keyring")
			}

			// Get Provenance File
			provenanceFile, err := chartManager.ReadProvenanceFile()

			if err != nil {
				return errors.Wrap(err, "Failed to Read Provenance File")
			}

			verifier := verifier.Verifier{
				ChartManager: chartManager,
				Keyring:      keyring,
			}

			err = verifier.VerifyChart()

			if err != nil {
				return err
			}

			uploadRequest := &rekor.RekorUploadRequest{
				Provenance: provenanceFile,
				PublicKey:  verifier.PublicKey,
			}

			uploadResponse, err := r.Upload(uploadRequest)

			if err != nil {
				return err
			}

			var newIndex int64
			for _, entry := range uploadResponse.Payload {
				newIndex = swag.Int64Value(entry.LogIndex)
			}

			fmt.Println(fmt.Sprintf("Created Helm entry at index %d, available at: %v%v\n", newIndex, rekorServer, string(uploadResponse.Location)))
			return nil
		},
	}

	addRekorFlags(uploadCmd, &uploadOptions)
	addPkiFlags(uploadCmd, &uploadOptions)

	return uploadCmd

}
