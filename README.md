[![Build Status](https://travis-ci.org/grafeas/voucher.svg?branch=master)](https://travis-ci.org/grafeas/voucher)
[![codecov](https://codecov.io/gh/grafeas/voucher/branch/master/graph/badge.svg)](https://codecov.io/gh/grafeas/voucher)
[![Go Report Card](https://goreportcard.com/badge/github.com/grafeas/voucher)](https://goreportcard.com/report/github.com/grafeas/voucher)

# voucher

**Table of Contents**

- [Introduction](#introduction)
- [Contributing](#contributing)

## Introduction

Voucher is the missing piece in the binary authorization toolchain which enables you to secure your software supply pipeline. Binary authorization uses an admission controller such as [Kritis](https://github.com/grafeas/kritis), which pulls information about a container image from a metadata server such as [Grafeas](https://grafeas.io/) to ensure that the image is not deployed to production unless it has passed an appropriate suite of checks. As running checks on containers during deployment is time consuming and prevents rapid rollout of changes, the checks the admission controller utilizes to verify an image is ready for production should be run at build time. Voucher does exactly that.

Voucher was designed to be called from your CI/DC pipeline, after an image is built, but before that image is deployed to production. Voucher pulls the newly built image from your image registry; runs it through all of the checks that were requested, and generates attestations for every check that the image passes. Those attestations (OpenPGP signatures of container digests) are then pushed to the metadata server, where Kritis can verify them.

Voucher presently includes the following checks:

| Test Name    | Description                                                                        |
| :--------    | :--------------------------------------------------------------------------------- |
| `diy`        | Can the image be downloaded from our container registry?                           |
| `nobody`     | Was the image built to run as a user who is not root?                              |
| `snakeoil`   | Is the image free of known security issues?                                        |
| `provenance` | Was the image built by us or a trusted system?                                     |
| `approved`   | Did the source code for the image pass all required checks in the code repository? |

As well as the following dynamic check:

| Test Name       | Description                                               |
| :-------------- | :-------------------------------------------------------- |
| `is_<org name>` | Did the source for this image come from the passed organization (for example, `is_shopify`) |

Note that `provenance` and the dynamic checks require the prescence of build metadata in your metadata store. While unsigned metadata is valid, to ensure that you are trusting metadata that hasn't been forged, it is recommended that you use signed metadata as well.

## Voucher Server and Client

This repository contains two tools: Voucher server, intended to run in your infrastructure to respond to CI/CD pipeline requests, and Voucher client, an example of a Voucher API client that you can use directly in your CI/CD pipeline or as a basis for your own code.

- [Server](cmd/voucher_server/README.md)
- [Client](cmd/voucher_client/README.md)

## Contributing

Please refer to the [Contributing document](CONTRIBUTING.md) if you are interested in contributing to voucher!
