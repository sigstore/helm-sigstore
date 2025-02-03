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

package pgp

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"io"
	"os"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/ProtonMail/go-crypto/openpgp/clearsign"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
	"github.com/pkg/errors"
)

func GetKeyring(keyRingPath string, publicKeyPath string) (openpgp.EntityList, error) {
	var keyringEntityList openpgp.EntityList
	var err error

	if publicKeyPath == "" {
		// Load Keyring
		keyringEntityList, err = loadKeyRing(keyRingPath)

		if err != nil {
			return nil, errors.Wrap(err, "Error Retrieving Keyring")
		}
	} else {
		keyringEntityList, err = getKeyRingFromPublicKey(publicKeyPath)

		if err != nil {
			return nil, errors.Wrap(err, "Error Retrieving Keyring from Public Key")
		}
	}

	return keyringEntityList, nil
}

func ExtractPublicKey(entity *openpgp.Entity) ([]byte, error) {
	gotWriter := bytes.NewBuffer(nil)
	wr, err := armor.Encode(gotWriter, openpgp.PublicKeyType, nil)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to Encode")
	}

	if entity.Serialize(wr) != nil {
		return nil, errors.Wrap(err, "Failed to Serialize")
	}

	if wr.Close() != nil {
		return nil, errors.Wrap(err, "Failed to Close")
	}

	return gotWriter.Bytes(), nil
}

func GetFingerprintFromPublicKey(content []byte) (string, error) {
	entitylist, err := openpgp.ReadArmoredKeyRing(bytes.NewBuffer(content))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(entitylist[0].PrimaryKey.Fingerprint), nil
}

func VerifySignature(file []byte, keyring openpgp.EntityList) (*openpgp.Entity, *io.Reader, error) {
	block, _ := clearsign.Decode(file)
	if block == nil {
		// There was no sig in the file.
		return nil, nil, errors.New("signature block not found")
	}

	var bufferRead bytes.Buffer
	armoredSignatureReader := io.TeeReader(block.ArmoredSignature.Body, &bufferRead)

	signer, err := openpgp.CheckDetachedSignature(
		keyring,
		bytes.NewBuffer(block.Bytes),
		block.ArmoredSignature.Body,
		&packet.Config{},
	)

	return signer, &armoredSignatureReader, err
}

func getKeyRingFromPublicKey(keypath string) (openpgp.EntityList, error) {
	f, err := os.Open(keypath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return openpgp.ReadArmoredKeyRing(bufio.NewReader(f))
}

func loadKeyRing(ringpath string) (openpgp.EntityList, error) {
	f, err := os.Open(ringpath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return openpgp.ReadKeyRing(f)
}
