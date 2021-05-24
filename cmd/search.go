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

	"github.com/pkg/errors"
	"github.com/sigstore/helm-sigstore/pkg/chart"
	"github.com/sigstore/helm-sigstore/pkg/rekor"
	"github.com/sigstore/helm-sigstore/pkg/types"
	"github.com/spf13/cobra"
)

func NewSearchCmd() *cobra.Command {

	searchOptions := types.CLIOptions{}

	// searchCmd represents the upload command
	searchCmd := &cobra.Command{
		Use:   "search",
		Short: "Search for a Signed Helm Chart",
		RunE: func(cmd *cobra.Command, args []string) error {

			if len(args) != 1 {
				return errors.New("1 argument (Path to packaged chart) is required")
			}

			chartPath := args[0]

			chartManager, err := chart.NewChartManager(chartPath)

			if err != nil {
				return err
			}

			rekor, err := rekor.NewRekor(searchOptions.RekorServer)

			if err != nil {
				return err
			}

			uuids, err := rekor.Search(chartManager)

			if err != nil {
				return err
			}

			if len(uuids) == 0 {
				return errors.New(fmt.Sprintf("Unable to find an entry for Chart '%s'\n", chartManager.ChartPath))
			}

			fmt.Println("The Following Records were Found")

			fmt.Println(fmt.Sprintf("\nRekor Server: %s", searchOptions.RekorServer))
			for _, uuid := range uuids {
				fmt.Println(fmt.Sprintf("Rekor UUID: %s", uuid))
			}

			return nil
		},
	}

	addRekorFlags(searchCmd, &searchOptions)

	return searchCmd
}
