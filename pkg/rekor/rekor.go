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

package rekor

import (
	"fmt"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/pkg/errors"
	"github.com/sigstore/helm-sigstore/pkg/chart"
	"github.com/sigstore/rekor/pkg/client"
	generatedclient "github.com/sigstore/rekor/pkg/generated/client"
	"github.com/sigstore/rekor/pkg/generated/client/entries"
	"github.com/sigstore/rekor/pkg/generated/client/index"
	"github.com/sigstore/rekor/pkg/generated/models"
	helm_v001 "github.com/sigstore/rekor/pkg/types/helm/v0.0.1"
)

type Rekor struct {
	rekorClient *generatedclient.Rekor
}

type RekorUploadRequest struct {
	Provenance []byte
	PublicKey  []byte
}

type RekorUploadResponse struct {
	Location strfmt.URI
	Payload  models.LogEntry
}

func NewRekor(rekorServer string) (*Rekor, error) {

	r, err := client.GetRekorClient(rekorServer)

	if err != nil {
		return nil, errors.Wrap(err, "Failed to Create Rekor Client")
	}

	return &Rekor{
		rekorClient: r,
	}, nil

}

func (r *Rekor) Search(chartManager *chart.ChartManager) ([]string, error) {

	hashVal, err := chartManager.GetChartDigest()

	if err != nil {
		return nil, errors.Wrap(err, "Failed to Chart Digest")
	}

	params := index.NewSearchIndexParams()
	params.Query = &models.SearchIndex{}
	params.Query.Hash = hashVal

	resp, err := r.rekorClient.Index.SearchIndex(params)

	if err != nil {
		return nil, errors.Wrap(err, "Error Querying Rekor Server")
	}

	return resp.GetPayload(), nil

}

func (r *Rekor) Upload(request *RekorUploadRequest) (*RekorUploadResponse, error) {

	params := entries.NewCreateLogEntryParams()

	re := new(helm_v001.V001Entry)
	re.HelmObj = models.HelmV001Schema{}
	re.HelmObj.Chart = &models.HelmV001SchemaChart{}
	re.HelmObj.Chart.Provenance = &models.HelmV001SchemaChartProvenance{}
	re.HelmObj.Chart.Provenance.Content = strfmt.Base64(request.Provenance)
	re.HelmObj.PublicKey = &models.HelmV001SchemaPublicKey{}
	re.HelmObj.PublicKey.Content = strfmt.Base64(request.PublicKey)

	entry := models.Helm{}
	entry.APIVersion = swag.String(re.APIVersion())
	entry.Spec = re.HelmObj

	params.SetProposedEntry(&entry)

	resp, err := r.rekorClient.Entries.CreateLogEntry(params)
	if err != nil {
		switch e := err.(type) {
		case *entries.CreateLogEntryConflict:
			return nil, errors.New(fmt.Sprintf("Entry already exists: %s", e.Location.String()))
		default:
			return nil, errors.Wrap(err, "Error Creating Log Entry")
		}
	}

	return &RekorUploadResponse{
		Location: resp.Location,
		Payload:  resp.Payload,
	}, nil

}

func (r *Rekor) GetByUUID(uuid string) (*models.LogEntryAnon, error) {

	params := entries.NewGetLogEntryByUUIDParams()
	params.EntryUUID = uuid

	resp, err := r.rekorClient.Entries.GetLogEntryByUUID(params)
	if err != nil {
		return nil, err
	}

	for k, entry := range resp.Payload {

		if k != uuid {
			continue
		}

		return &entry, nil

	}

	return nil, errors.New("Unable to find LogEntry matching UUID")

}
