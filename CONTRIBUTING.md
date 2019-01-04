# How to contribute

Thank you for considering contributing to Voucher!

If you are interested in adding a new check to Voucher's code base, please read
the [Checks documentation](/checks/README.md) first, to get an idea as to how
Checks work.

## Getting started

- Review this document and the [Code of Conduct](CODE_OF_CONDUCT.md).

- Setup a [Go development environment](https://golang.org/doc/install#install)
if you haven't already.

- Get the Shopify version the project by using Go get:

```
$ go get -u github.com/Shopify/voucher
```

- Fork this project on GitHub.

- Setup your fork as a remote for your project:

```
$ cd $GOPATH/src/github.com/Shopify/voucher
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

    - If you update the dependencies, you may have to run `make update-deps` to
 ensure that the dependencies are vendored. This should happen in a separate
commit from those containing your source modifications.

- Push your changes to your fork's remote:

```
$ git push -u <your username> <the name of your branch>
```

## Send in your changes

- Sign the [Contributor License Agreement](https://cla.shopify.com).

- Run the test suite and make certain that the tests are not failing. If you
are adding code which would be untested, please consider adding tests to cover
that code.

- Open a PR against Shopify/voucher
