# How to contribute

Thank you for considering contributing to Voucher!

If you are interested in adding a new check to Voucher's code base, please read
the [Checks documentation](/checks/README.md) first, to get an idea as to how
Checks work.

## Getting started

- Review this document and the [Code of Conduct](CODE_OF_CONDUCT.md).

- Setup a [Go development environment](https://golang.org/doc/install#install)
if you haven't already.

- Get the source code by using Go get:

```
$ go get -u github.com/grafeas/voucher
```

- Fork this project on GitHub.

- Setup your fork as a remote for your project:

```
$ cd $GOPATH/src/github.com/grafeas/voucher
$ git remote add <your username> <your fork's remote path>
```

## Work on your feature

- Create your feature branch based off of the `master` branch. (It might be
worth doing a `git pull` if you haven't done one in a while.)

```
$ git checkout master
$ git pull
$ git checkout -b <the name of your branch>
```

- Code!

    - Please run `go fmt` and `golint` while you work on your change, to clean
up your formatting/check for issues.

    - If you add a dependencies, run `go mod tidy` to prune any no-longer-needed
    dependencies from `go.mod`.

    - To update dependencies to use newer minor or patch releases when available:
    run `make update-deps`.

- Push your changes to your fork's remote:

```
$ git push -u <your username> <the name of your branch>
```

## Send in your changes

- Sign the [Contributor License Agreement](https://cla.developers.google.com/)

- Run the test suite and make certain that the tests are not failing. If you
are adding code which would be untested, please consider adding tests to cover
that code.

- Open a PR against grafeas/voucher

## Making a release

If you are maintaining the Voucher project you can make a release of Voucher
using `goreleaser`.

First, update the master branch to the commit that you'd like to create a
release for.

```shell
$ git checkout master
$ git pull
```

If everything looks good, create a new tag. Voucher uses
[Semantic Versioning](https://semver.org) which basically means that API
compatible changes should bump the last digit, backwards compatible changes
should bump the second, and API incompatible changes should bump the first
digit. For example, if we are fixing a bug that doesn't affect other systems
we can bump from v1.0.0 to v1.0.1. If the change is backwards compatible with
the previous version, we can bump from v1.0.0 to v1.1.0. If the change is
not backwards compatible, we must bump the version from v1.0.0 to v2.0.0.

Run the following, where `version` is replaced with the appropriate version for
this release.

(Note that you will need to have Git configured to sign tags with 
your OpenPGP key or this command will fail.)

```shell
$ git tag -s <version>
```

Before pushing your tag, build a release version of Voucher to ensure that
everything builds properly:

```shell
$ make release
```

If this step fails, please do not make a release without fixing it. In
addition, please delete the tag so it can be replaced once the issue is
fixed.

If this step is successful, the ready-to-be released files will be found
under `dist/`.

Next, push the tag to the server, where `version` is the same version you
specified before:

```shell
$ git push origin refs/tags/<version>
```

Finally, create a new release in GitHub for the new version for the tag you
created and signed. In the `dist/` directory you will find the automatically
generated binary tar archives and `checksums.txt` file, which will need to
be added to the release in GitHub.

In the release description, paste the output of the `CHANGELOG.md` file,
also automatically generated in `dist/`. This should be touched up to
remove unnecessary commit descriptions.

Once you're ready, publish the release and tell everyone about it!
