# Voucher Subscriber

- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)

## Installation

Install `voucher_subscriber` by running:

```shell
$ go get -u github.com/grafeas/voucher/cmd/voucher_subscriber
```

This will download and install the voucher subscriber binary into `$GOPATH/bin` directory.

## Configuration

See the [Tutorial](TUTORIAL.md) for more thorough setup instructions.

An example configuration file can be found in the [config directory](../../config/config.toml).

The configuration can be written as a toml, json, or yaml file, and you can specify the path to the configuration file using the `-c` flag.

All of the configuration options for the Voucher Subscriber is the same as the [Voucher Server](../voucher_server/README.md#configuration)

## Usage

You can run Voucher in pub/sub subscriber mode by launching `voucher_subscriber`, using the following syntax:

```shell
$ voucher_subscriber [--project <project> --subscription <subscription>]
```

`voucher_subscriber` supports the following flags:

| Flag        | Short Flag       | Description                                                                |
| :--------   | :--------------- | :------------------------------------------------------------------------- |
| `--project`  | `-p`            | The GCP project to be used.                                                |
| `--subscription` | `-s`        | The subscription that contains messages.                                   |
| `--timeout` |                  | The number of seconds to spend checking an image, before failing.          |

For example:

```shell
$ voucher_subscriber --project my-project --subscription my-subscription
```

This would launch the subscriber listening to `my-subscription` in `my-project`.

The subscriber expects messages with the following format:

```json
{
  "action":"INSERT",
  "digest":"gcr.io/my-project/hello-world@sha256:6ec128e26cd5...",
  "tag":"gcr.io/my-project/hello-world:1.1"
}
```
The `tag` field can be omitted but the `action` and `digest` fields are required.

The subscriber currently retries when it can't connect to a MetaData client or start performing a `CheckSuite`. It will log whenever this does happen.

More details about configuring pub/sub with GCR can be found in the [official documentation](https://cloud.google.com/container-registry/docs/configuring-notifications).
