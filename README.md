[![Build Status](https://travis-ci.org/Shopify/voucher.svg?branch=master)](https://travis-ci.org/Shopify/voucher)
[![codecov](https://codecov.io/gh/Shopify/voucher/branch/master/graph/badge.svg)](https://codecov.io/gh/Shopify/voucher)
[![Go Report Card](https://goreportcard.com/badge/github.com/Shopify/voucher)](https://goreportcard.com/report/github.com/Shopify/voucher)

# voucher

**Table of Contents**

- [Introduction](#introduction)
- [Installing voucher](#installing-voucher)
  - [Voucher Server](#voucher-server)
  - [Voucher Standalone](#voucher-standalone)
  - [Voucher Client](#voucher-client)
- [Voucher Server and Standalone](#voucher-server-and-standalone)
  - [Configuration](#configuration)
    - [Scanner](#scanner)
    - [Fail-On: Failing on vulnerabilities](#fail-on-failing-on-vulnerabilities)
    - [Valid Repos](#valid-repos)
    - [Enabling Checks](#enabling-checks)
  - [Running Voucher](#running-voucher)
    - [Using voucher standalone to check an image](#using-voucher-standalone-to-check-an-image)
    - [Using Voucher Server to check an image](#using-voucher-server-to-check-an-image)
- [Voucher Client](#voucher-client)
  - [Configuration](#configuration)
  - [Using Voucher Client](#using-voucher-client)
- [Contributing](#contributing)

## Introduction

Voucher is the missing piece in the binary authorization toolchain which enables you to secure your software supply pipeline. Binary authorization uses an admission controller such as [Kritis](https://grafeas.io/docs/concepts/what-is-kritis/overview.html), which pulls information about a container image from a metadata server such as [Grafeas](https://grafeas.io/) to ensure that the image is not deployed to production unless it has passed an appropriate suite of checks. As running checks on containers during deployment is time consuming and prevents rapid rollout of changes, the checks the admission controller utilizes to verify an image is ready for production should be run at build time. Voucher does exactly that.

Voucher was designed to be called from your CI/DC pipeline, after an image is built, but before that image is deployed to production. Voucher pulls the newly built image from your image registry; runs it through all of the checks that were requested, and generates attestations for every check that the image passes. Those attestations (OpenPGP signatures of container digests) are then pushed to the metadata server, where Kritis can verify them.

Voucher presently includes the following checks:

| Test Name    | Description                                                                    |
| :--------    | :----------------------------------------------------------------------------- |
| `diy`        | Can the image be downloaded from our container registry?                       |
| `nobody`     | Was the image built to run as a user who is not root?                          |
| `snakeoil`   | Is the image free of known security issues?                                    |
| `provenance` | Was the image built by us or a trusted system?                                 |

## Installing voucher

### Voucher Server

Install `voucher_server` by running:

```shell
$ go get -u github.com/Shopify/voucher/cmd/voucher_server
```

This will download and install the voucher server binary into `$GOPATH/bin` directory.

### Voucher Standalone

Install the standalone version of Voucher, `voucher_cli`, by running:

```shell
$ go get -u github.com/Shopify/voucher/cmd/voucher_cli
```

This will download and install the voucher binary into `$GOPATH/bin` directory.

### Voucher Client

Voucher Client is a tool for connecting to a running Voucher server.

Install the client, `voucher_client`, by running:

```shell
$ go get -u github.com/Shopify/voucher/cmd/voucher_client
```

This will download and install the voucher client into `$GOPATH/bin` directory.

## Voucher Server and Standalone

See the [Tutorial](TUTORIAL.md) for more thorough setup instructions.

### Configuration

An example configuration file can be found in the [config directory](config/config.toml).

The configuration can be written as a toml, json, or yaml file, and you can specify the path to the configuration file using "-c".

Below are the configuration options for Voucher Standalone and Server:

| Group     | Key               | Description                                                                                           |
| :-------- | :---------------  | :---------------------------------------------------------------------------------------------------- |
|           | `dryrun`          | When set, don't create attestations.                                                                  |
|           | `scanner`         | The vulnerability scanner to use ("clair" or "gca").                                                  |
|           | `failon`          | The minimum vulnerability to fail on. Discussed below.                                                |
|           | `valid_repos`     | A list of repos that are owned by your team/organization.                                             |
|           | `image_project`   | The project in the metadata server that image information is stored.                                  |
|           | `binauth_project` | The project in the metadata server that the binauth information is stored.                            |
|           | `timeout`         | The number of seconds to spend checking an image, before failing (voucher standalone only).           |
| `checks`  | (test name here)  | A test that is active when running "all" tests.                                                       |
| `server`  | `port`            | The port that the server can be reached on.                                                           |
| `server`  | `timeout`         | The number of seconds to spend checking an image, before failing.                                     |
| `server`  | `require_auth`    | Require the use of Basic Auth, with the username and password from the configuration.                 |
| `server`  | `username`        | The username that Voucher server users must use.                                                      |
| `server`  | `password`        | A password hashed with the bcrypt algorithm, for use with the username.                               |
| `ejson`   | `dir`             | The path to the ejson keys directory.                                                                 |
| `ejson`   | `secrets`         | The path to the ejson secrets.                                                                        |
| `clair`   | `address`         | The hostname that Clair exists at. If "http://" or "https://" is omitted, this will default to HTTPS. |

Configuration options can be overridden at runtime by setting the appropriate flag. For example, if you set the "port" flag when running `voucher_server`, that value will override whatever is in the configuration.

#### Scanner

The `scanner` option in the configuration is used to select the Vulnerability scanner.

This option supports two values:

- `c` or `clair` to use an instance of CoreOS's Clair.
- `g` or `gca` to use Google Container Analysis.

If you decide to use Clair, you will need to update the clair configuration block to specify the correct address for the server.

#### Fail-On: Failing on vulnerabilities

The `failon` option allows you to set the minimum vulnerability to consider an image insecure.

This option supports the following:

- "negligible"
- "low"
- "medium"
- "unknown"
- "high"
- "critical"

For example, if you set `failon` to "high", only "high" and "critical" vulnerabilities will prevent the image from being attested. A value of "low" will cause "low", "medium", "unknown", "high", and "critical" vulnerabilities to prevent the image from being attested failure.

#### Valid Repos

The `valid_repos` option in the configuration is used to limit which repositories images must be from to pass the DIY check.

This option takes a list of repos, which are compared against the repos that images live in. An image will pass if it starts with any of the items in the list.

For example:

```json
{
    "valid_repos": [
        "gcr.io/team-images/",
        "gcr.io/external-images/specific-project",
    ]
}
```

Will allow images that start with `gcr.io/team-images/` and `gcr.io/external-images/specific-project/` to pass the DIY check, while blocking other `gcr.io/external-images/`.

#### Enabling Checks

You can enable certain checks for the "all" option by updating the `checks` block in the configuration.

For example:

```toml
[checks]
diy      = true
nobody   = true
provenance = false
snakeoil = true
```

With this configuration, the "diy", "nobody", and "snakeoil" checks would run when running "all" checks. The "provenance" check will be ignored unless called directly.

### Running Voucher

#### Using voucher standalone to check an image

You can run Voucher's standalone version by `voucher_cli`, using the following syntax:

```shell
$ voucher_cli <test to run> --image <image to check> [other options]
```

`voucher_cli` supports the following flags:

| Flag        | Short Flag       | Description                                                                |
| :--------   | :--------------- | :------------------------------------------------------------------------- |
| `--config`  | `-c`             | The path to a configuration file that should be used.                      |
| `--dryrun`  |                  | When set, don't create attestations.                                       |
| `--scanner` |                  | The vulnerability scanner to use ("clair" or "gca").                       |
| `--failon`  |                  | The minimum vulnerability to fail on. Discussed above.                     |
| `--image`   | `-i`             | The image to check and attest.                                             |
| `--timeout` |                  | The number of seconds to spend checking an image, before failing.          |

For example:

```shell
$ voucher_cli all --image gcr.io/path/to/image@sha256:ab7524b7375fbf09b3784f0bbd9cb2505700dd05e03ce5f5e6d262bf2f5ac51c
```

This would run all possible tests, or all tests that are enabled by the [configuration](#configuration), against the image located at the passed URL.

You can also run an individual test, by specifying that test:

```shell
$ voucher_cli diy --image gcr.io/path/to/image@sha256:ab7524b7375fbf09b3784f0bbd9cb2505700dd05e03ce5f5e6d262bf2f5ac51c
```

This would run the "diy" test.

#### Using Voucher Server to check an image

You can run Voucher in server mode by launching `voucher_server`, using the following syntax:

```shell
$ voucher_server [--port <port number>]
```

`voucher_server` supports the following flags:

| Flag        | Short Flag       | Description                                                                |
| :--------   | :--------------- | :------------------------------------------------------------------------- |
| `--config`  | `-c`             | The path to a configuration file that should be used.                      |
| `--port`    | `-p`             | Set the port to listen on.                                                 |
| `--timeout` |                  | The number of seconds to spend checking an image, before failing.          |

For example:

```shell
$ voucher_server --port 8000
```

This would launch the server, utilizing port 8000.

You can connect to Voucher over http.

For example, using `curl`:

```shell
$ curl -X POST -H "Content-Type: application/json" -d "{\"image_url\": \"gcr.io/path/to/image@sha256:ab7524b7375fbf09b3784f0bbd9cb2505700dd05e03ce5f5e6d262bf2f5ac51c\"}" http://localhost:8000/all
```

The response will be something along the following lines:

```json
{
    "image": "gcr.io/path/to/image@sha256:ab7524b7375fbf09b3784f0bbd9cb2505700dd05e03ce5f5e6d262bf2f5ac51c",
    "success": false,
    "results": [
        {
            "name": "provenance",
            "error": "no occurrences returned for image",
            "success": false,
            "attested": false
        },
        {
            "name": "snakeoil",
            "success": true,
            "attested": true
        },
        {
            "name": "diy",
            "success": true,
            "attested": true
        },
        {
            "name": "nobody",
            "success": true,
            "attested": true
        }
    ]
}
```

More details about Voucher server can be read in the [API documentation](server/README.md).

## Voucher Client

### Configuration

The configuration for Voucher Client can be written as a toml, json, or yaml file, and you can specify the path to the configuration file using "-c". By default, the configuration is expected to be located at `~/.voucher{.yaml,.toml,.json}`.

Below are the configuration options for Voucher Standalone and Server:

| Key         | Description                                                                                |
| :---------- | :----------------------------------------------------------------------------------------- |
| `hostname`  | The Voucher server to connect to.                                                          |
| `timeout`   | The number of seconds to wait before failing (defaults to 240).                            |
| `username`  | Username to authenticate against Voucher with.                                             |
| `password`  | Password to authenticate against Voucher with.                                             |

Configuration options can be overridden at runtime by setting the appropriate flag. For example, if you set the "port" flag when running `voucher_server`, that value will override whatever is in the configuration.

 Here is an example (yaml encoded) configuration file:

```yaml
---
hostname: "https://my-voucher-server"
username: "<username>"
password: "<password>"
```

### Using Voucher Client

While you can use `curl` to make API calls against Voucher, you can also use `voucher_client` to save from making HTTP requests by hand. Unlike the other Voucher tools, `voucher_client` will look up the appropriate canonical version of an image reference if passed a tagged image reference.

```shell
$ voucher_client [--voucher <server> --check <check to run>] <image path>
```

`voucher_client` supports the following flags:

| Flag        | Short Flag       | Description                                                                |
| :--------   | :--------------- | :------------------------------------------------------------------------- |
| `--config`  |                  | The path to your configuration file, (default is $HOME/.voucher.yaml)      |
| `--check`   | `-c`             | The Check to run on the image ("all" for all checks).                      |
| `--voucher` | `-v`             | The Voucher server to connect to.                                          |
| `--username`|                  | Username to authenticate against Voucher with.                             |
| `--password`|                  | Password to authenticate against Voucher with.                             |
| `--timeout` | `-t`             | The number of seconds to wait before failing (defaults to 240).            |

For example:

```shell
$ voucher_client -v http://localhost:8000 gcr.io/path/to/image:latest
```

The output will be something along the following lines:

```json
 - Attesting image: gcr.io/path/to/image@sha256:ab7524b7375fbf09b3784f0bbd9cb2505700dd05e03ce5f5e6d262bf2f5ac51c
   ✗ nobody failed
   ✗ snakeoil failed, err: vulnernable to 1 vulnerabilities: CVE-2018-12345 (high)
   ✓ diy succeeded, but wasn't attested, err: rpc error: code = AlreadyExists desc = Requested entity already exists
```

## Contributing

Please refer to the [Contributing document](CONTRIBUTING.md) if you are interested in contributing to voucher!
