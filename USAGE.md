# Detailed Usage

## General Options

**Rekor Server** 

An instance of [Rekor](https://github.com/sigstore/rekor) must be available since the majority of the commands interact with an instance. By default, the public instance ([https://rekor.sigstore.dev](https://rekor.sigstore.dev)) is used. An alternate instance can be specified by using the `--rekor-server` flag or setting the `REKOR_SERVER` environment variable.

**GPG**

The public key that can be used to validate the signature of the Helm provenance file is required. The key can be standalone or contained within a Keyring. By default, a keyring located at `~/.gnupg/pubring.gpg` is used. An alternate location can be provided by specifying the `--keyring` flag or setting the `KEYRING` environment variable.

## Create, Package and Sign a Helm Chart and Upload to Rekor

Assuming you have the necessary tools to [sign a Helm Chart](https://helm.sh/docs/topics/provenance/), let's demonstrate the steps to create, package, sign and upload the newly created chart to Rekor.

First, create a new chart called `my-sigstore-chart`

```shell
$ helm create my-sigstore-chart

Creating my-sigstore-chart
```

Package and sign the chart

```shell
$ helm package --sign --key="<key>" my-sigstore-chart/

Successfully packaged chart and saved it to: <base_directory>/my-sigstore-chart-0.1.0.tgz
```

Upload the Chart to Rekor

```shell
$ helm sigstore upload <base_directory>/my-sigstore-chart-0.1.0.tgz

Created Helm entry at index 6846, available at: https://rekor.sigstore.dev/api/v1/log/entries/15b6a8215057f47f96ef34b58c5537e6a8eb8e2f83c2fcd8a831f9d2e813be0d
```

With the entry added to Rekor, it can be used to demonstrate the other subcommands in the subsequent sections.

## Searching for Signed Helm Charts

You can determine whether a Signed Helm Chart has an existing entry in the Rekor server using the `search` subcommand. Using the signed Helm chart created in the prior section, search or the entry in the Rekor server:

```shell
$ helm sigstore search <base_directory>/my-sigstore-chart-0.1.0.tgz

The Following Records were Found

Rekor Server: https://rekor.sigstore.dev
Rekor UUID: 15b6a8215057f47f96ef34b58c5537e6a8eb8e2f83c2fcd8a831f9d2e813be0d
```

## Verifying Signed Helm Chart

Signed Helm Charts can be verified based on entries in Rekor. The verification process will first search for an entry in Rekor, and if found, retrieve the entry and compare characteristics, such as public key, Chart hash and signature between the Rekor entry and the Signed chart.

Verify the signed Helm chart created in the prior sections:

```shell
$ helm sigstore verify <base_directory>/my-sigstore-chart-0.1.0.tgz

Chart Verified Successfully From Helm entry:

Rekor Server: https://rekor.sigstore.dev
Rekor Index: 6846
Rekor UUID: 15b6a8215057f47f96ef34b58c5537e6a8eb8e2f83c2fcd8a831f9d2e813be0d
```
