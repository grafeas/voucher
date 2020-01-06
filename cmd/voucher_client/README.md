# Voucher Client

## Installation

Voucher Client is a tool for connecting to a running Voucher server.

Install the client, `voucher_client`, by running:

```shell
$ go get -u github.com/Shopify/voucher/cmd/voucher_client
```

This will download and install the voucher client into `$GOPATH/bin` directory.

## Configuration

The configuration for Voucher Client can be written as a toml, json, or yaml file, and you can specify the path to the configuration file using "-c". By default, the configuration is expected to be located at `~/.voucher{.yaml,.toml,.json}`.

Below are the configuration options for Voucher Client:

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

## Using Voucher Client

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
