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

package chart

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"strings"

	"github.com/pkg/errors"
)

type ChartManager struct {
	ChartPath           string
	ChartProvenancePath string
	chartDigest         string
}

func NewChartManager(chartPath string) (*ChartManager, error) {

	_, err := os.Stat(chartPath)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to Load Chart")
	}

	if !strings.EqualFold(filepath.Ext(chartPath), ".tgz") {
		return nil, errors.New("Chart is not a .tgz file")
	}

	provfile := chartPath + ".prov"
	if _, err := os.Stat(provfile); err != nil {
		return nil, errors.Wrapf(err, "could not load provenance file %s", provfile)
	}

	return &ChartManager{
		ChartPath:           chartPath,
		ChartProvenancePath: provfile,
	}, nil

}

func (c *ChartManager) ReadProvenanceFile() ([]byte, error) {
	return readFile(c.ChartProvenancePath)
}

func readFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func (c *ChartManager) GetChartDigest() (string, error) {

	if c.chartDigest != "" {
		return c.chartDigest, nil
	}

	file, err := os.Open(filepath.Clean(c.ChartPath))

	if err != nil {
		return "", fmt.Errorf("error opening chart '%v': %w", c.ChartPath, err)
	}

	defer file.Close()

	hasher := sha256.New()
	tee := io.TeeReader(file, hasher)

	if _, err := ioutil.ReadAll(tee); err != nil {
		return "", fmt.Errorf("error processing '%v': %w", c.ChartPath, err)
	}

	c.chartDigest = strings.ToLower(hex.EncodeToString(hasher.Sum(nil)))

	return c.chartDigest, nil

}
