# Voucher Server

- [Installation](#installation)
- [Configuration](#configuration)
  - [Scanner](#scanner)
  - [Fail-On: Failing on vulnerabilities](#fail-on-failing-on-vulnerabilities)
  - [Valid Repos](#valid-repos)
  - [Trusted Builder Identities and Trusted Builder Projects](#trusted-builder-identities-and-trusted-builder-projects)
  - [Repository Checks](#repository-checks)
    - [Repository Groups](#repository-groups)
    - [Organization Check](#organization-check)
  - [Enabling Checks](#enabling-checks)
- [Usage](#usage)

## Installation

Install `voucher_server` by running:

```shell
$ go get -u github.com/Shopify/voucher/cmd/voucher_server
```

This will download and install the voucher server binary into `$GOPATH/bin` directory.

## Configuration

See the [Tutorial](TUTORIAL.md) for more thorough setup instructions.

An example configuration file can be found in the [config directory](../../config/config.toml).

The configuration can be written as a toml, json, or yaml file, and you can specify the path to the configuration file using "-c".

Below are the configuration options for Voucher Server:

| Group                | Key                          | Description                                                                                           |
| :-------------       | :--------------------------- | :---------------------------------------------------------------------------------------------------- |
|                      | `dryrun`                     | When set, don't create attestations.                                                                  |
|                      | `scanner`                    | The vulnerability scanner to use ("clair" or "gca").                                                  |
|                      | `failon`                     | The minimum vulnerability to fail on. Discussed below.                                                |
|                      | `valid_repos`                | A list of repos that are owned by your team/organization.                                             |
|                      | `trusted_builder_identities` | A list of email addresses. Owners of these emails are considered "trusted" (and will pass Provenance) |
|                      | `trusted_projects`           | A list of projects that are considered "trusted" (and will pass Provenance)                           |
|                      | `image_project`              | The project in the metadata server that image information is stored.                                  |
|                      | `binauth_project`            | The project in the metadata server that the binauth information is stored.                            |
| `checks`             | (test name here)             | A test that is active when running "all" tests.                                                       |
| `server`             | `port`                       | The port that the server can be reached on.                                                           |
| `server`             | `timeout`                    | The number of seconds to spend checking an image, before failing.                                     |
| `server`             | `require_auth`               | Require the use of Basic Auth, with the username and password from the configuration.                 |
| `server`             | `username`                   | The username that Voucher server users must use.                                                      |
| `server`             | `password`                   | A password hashed with the bcrypt algorithm, for use with the username.                               |
| `ejson`              | `dir`                        | The path to the ejson keys directory.                                                                 |
| `ejson`              | `secrets`                    | The path to the ejson secrets.                                                                        |
| `clair`              |  `address`                   | The hostname that Clair exists at. If "http://" or "https://" is omitted, this will default to HTTPS. |
| `repository.[alias]` | `org-url`                    | The URL used to determine if a repository is owned by an organization.                                |
| `required.[env]`     | (test name here)             | A test that is active when running "env" tests.                                                       |

Configuration options can be overridden at runtime by setting the appropriate flag. For example, if you set the "port" flag when running `voucher_server`, that value will override whatever is in the configuration.

Note that Repositories can be set multiple times. This is discussed futher below.

### Scanner

The `scanner` option in the configuration is used to select the Vulnerability scanner.

This option supports two values:

- `c` or `clair` to use an instance of CoreOS's Clair.
- `g` or `gca` to use Google Container Analysis.

If you decide to use Clair, you will need to update the clair configuration block to specify the correct address for the server.

### Fail-On: Failing on vulnerabilities

The `failon` option allows you to set the minimum vulnerability to consider an image insecure.

This option supports the following:

- "negligible"
- "low"
- "medium"
- "unknown"
- "high"
- "critical"

For example, if you set `failon` to "high", only "high" and "critical" vulnerabilities will prevent the image from being attested. A value of "low" will cause "low", "medium", "unknown", "high", and "critical" vulnerabilities to prevent the image from being attested failure.

### Valid Repos

The `valid_repos` option in the configuration is used to limit which repositories images must be from to pass the DIY check.

This option takes a list of repos, which are compared against the repos that images live in. An image will pass if it starts with any of the items in the list.

For example:

```toml
valid_repos = [
    "gcr.io/team-images/",
    "gcr.io/external-images/specific-project"
]
```

Will allow images that start with `gcr.io/team-images/` and `gcr.io/external-images/specific-project/` to pass the DIY check, while blocking other `gcr.io/external-images/`.

### Trusted Builder Identities and Trusted Builder Projects

The provenance check works by obtaining the build metadata for an image from the metadata service, and verifying that it both comes from a trusted project and was built by a trusted builder.

You can use this to ensure that only images built by your continuous integration pipeline can be deployed to the cloud, with exceptions for images built by trusted administrators.

For example:

```toml
trusted_builder_identities = [
    "catherine@example.com",
    "idcloudbuild.gserviceaccount.com"
]
trusted_projects = [
    "team-images"
]
```

This would require that images be built in the `team-images` project, by either `catherine@example.com` or `idcloudbuild.gserviceaccount.com`. Images not built in the `team-images` project will fail Provenance regardless of who built them.

### Repository Checks

#### Repository Groups

Repository Groups are used to determine which Repository client should be used
to connect to a Repository.

They are defined as an alias (usually matching the name of the organization in the repository system) and a URL. Note that
the alias can contain only lower cases letters, dashes and underscores (`[a-z_-]`)

The URL is used to determine if a repository is owned by a repository group.

For example, a repository owned by `Shopify`, it's URL should contain `github.com/Shopify`. A repository owned by `Grafeas`, should contain `github.com/grafeas`, and so on.

Repository groups are required to use the Repository checks.

You can define repository groups as follows:

```toml
[repository.shopify]
org-url = "https://github.com/Shopify"

[repository.grafeas]
org-url = "https://github.com/grafeas"
```

#### Organization Check

The organization check is a dynamic check which uses the name of an organization to determine if code came from that organization.

For example, if you have defined an organization as follows:

```toml
[[repositories]]
alias = "Shopify"
org-url = "https://github.com/Shopify"
```

You can enable a check that verifies that an image came from that organization:

```toml
[checks]
is_shopify = true
```

The name of the check is `is_` followed by the name of the organization, converted to lowercase.

This check works by:

- looking up the build metadata for an image from your metadata service
- getting the repository information from that metadata
- connecting to the API of that repository (in this example, Github)
- verifying that the source code is associated with the organization that it says it is

### Enabling Checks

You can enable certain checks for the "all" option by updating the `checks` block in the configuration.

For example:

```toml
[checks]
diy      = true
nobody   = true
provenance = false
snakeoil = true
is_shopify = true
```

With this configuration, the `diy`, `nobody`, `snakeoil`, and `is_shopify` checks would run when running `all` checks. The `provenance` check will be ignored unless called directly.

### Required Checks

You can configure named groups of checks the same as the `checks` block except by replacing `checks` with `required.[env]` where `[env]` is a name of your choosing.

For example:

```toml
[required.myenv]
diy      = true
nobody   = true
provenance = false
snakeoil = true
is_shopify = true
```

With this configuration, the `diy`, `nobody`, `snakeoil`, and `is_shopify` checks would run when running `myenv` checks. The `provenance` check will be ignored unless called directly.


## Usage

### Using Voucher Server to check an image

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

More details about Voucher server can be read in the [API documentation](../../server/README.md).


