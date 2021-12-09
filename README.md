# helm-sigstore

[![Build Status](https://github.com/sigstore/helm-sigstore/workflows/CI/badge.svg?branch=main)](https://github.com/sigstore/helm-sigstore/actions?workflow=CI)

Plugin for [Helm](https://helm.sh/) to integrate the [sigstore](https://sigstore.dev/) ecosystem. Search, upload and verify signed Helm Charts in the [Rekor](https://github.com/sigstore/rekor) Transparency Log. 

## Info

helm-sigstore is developed as part of the [`sigstore`](https://sigstore.dev) project.

We also use a [slack channel](https://sigstore.slack.com)!
Click [here](https://join.slack.com/t/sigstore/shared_invite/zt-mhs55zh0-XmY3bcfWn4XEyMqUUutbUQ) for the invite link.

## Installation

Use the following steps to build the `helm-sigstore` binary and install it as a Helm Plugin

### Building

On a system with [Go](https://golang.org/) installed, execute the following to download the source and build the plugin

```shell
$ mkdir -p $GOPATH/src/github.com/sigstore
$ cd $GOPATH/src/github.com/sigstore
$ git clone https://github.com/sigstore/helm-sigstore.git
$ cd helm-sigstore
```

Build the plugin

```shell
$ make
```

The plugin binary will be available in the `bin` directory

### Plugin Installation

Before installing `helm-sigstore` as a Helm plugin, ensure that Helm is installed and configured on your machine. Then install the plugin.

```shell
$ helm plugin install https://github.com/sigstore/helm-sigstore
```

Confirm the plugin is available in Helm

```
$ helm plugin list

NAME            VERSION         DESCRIPTION                                                                  
sigstore        0.1.0           This plugin integrates Helm into the Sigstore ecosystem.                     
```

With the installation complete and successful, the plugin can be invoked through the `helm sigstore` command

```shell
$ helm sigstore

Integrates sigstore with Helm

Usage:
  sigstore [command]
...
```

## Quickstart

This brief example demonstrates how to upload a signed Helm chart to Rekor and validate the entry

### Upload a Signed Helm Chart

```
$ helm sigstore upload <path_to_packaged_chart>

Created Helm entry at index 6821, available at: https://rekor.sigstore.dev/api/v1/log/entries/b30a142ef6c8b0480cd3e081fc99bc3d2a1a50ef60f68749c983a1479be6c4b9
```

_NOTE_: The provenance file must be located in the same directory as the packaged chart.
> To generate a provenance file, please consult the official documentation of [Helm Provenance and Integrity](https://helm.sh/docs/topics/provenance/).

### Verify the Signed Chart from Rekor

Use the same signed Helm chart from the prior section to verify the entry in Rekor

```shell
helm sigstore verify <path_to_packaged_chart>
Chart Verified Successfully From Helm entry:

Rekor Server: https://rekor.sigstore.dev
Rekor Index: 6821
Rekor UUID: b30a142ef6c8b0480cd3e081fc99bc3d2a1a50ef60f68749c983a1479be6c4b9
```

See the [Usage documentation](USAGE.md) for detailed explanations and additional options. 

## Security

Should you discover any security issues, please refer to sigstores [security
process](https://github.com/sigstore/community/blob/main/SECURITY.md)

