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
	"os"
	"path/filepath"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/sigstore/helm-sigstore/pkg/constants"
	"github.com/sigstore/helm-sigstore/pkg/types"
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "sigstore",
		Short: "Integrates sigstore with Helm",
		Long:  "Integrates sigstore with Helm",
	}

	rootCmd.SilenceErrors = true
	rootCmd.SilenceUsage = true

	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		return NewCLIFlagError(cmd, err)
	})

	rootCmd.AddCommand(NewSearchCmd())
	rootCmd.AddCommand(NewUploadCmd())
	rootCmd.AddCommand(NewVerifyCmd())
	rootCmd.AddCommand(NewVersionCmd())

	return rootCmd

}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	rootCmd := NewRootCommand()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func addRekorFlags(cmd *cobra.Command, options *types.CLIOptions) {

	cmd.PersistentFlags().StringVar(&options.RekorServer, "rekor-server", getEnv(constants.REKOR_SERVER_VAR, constants.DEFAULT_REKOR_SERVER), "server address:port")

}

func addPkiFlags(cmd *cobra.Command, options *types.CLIOptions) {

	cmd.PersistentFlags().StringVar(&options.PublicKey, "public-key", "", "location of the public key used to sign the chart")
	cmd.PersistentFlags().StringVar(&options.Keyring, "keyring", getDefaultKeyring(), "location of a public keyring")

}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getDefaultKeyring() string {

	if env, ok := os.LookupEnv("KEYRING"); ok {
		return env
	}

	if env, ok := os.LookupEnv("GNUPGHOME"); ok {
		return filepath.Join(env, "pubring.gpg")
	}

	// Locate Home Directory
	homedir, err := homedir.Dir()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return filepath.Join(homedir, ".gnupg", "pubring.gpg")
}
