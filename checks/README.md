# Voucher Checks

Adding new checks follows the same pattern that adding new SQL drivers uses. Checks are located in the same directory as this README, `/checks`, in their own package. While the Check's package does not need to match the name of the check, it helps with organization if it does.

The check name itself, which is set when registering the Check for use, will also be used to enable/disable that check, and to co-ordinate which PGP key to use when attesting images.

Below is an example of a Check, "examplecheck", which we will step through. This Check would logically be found at `checks/examplecheck/check.go`.

```golang
package examplecheck

import (
	"github.com/grafeas/voucher"
)

// check is a voucher.Check that holds our examplecheck test.
type check struct {
}

// Check if the image described by ImageData is good enough for our purposes.
func (c *check) Check(i voucher.ImageData) (bool, error) {
    ok := isImageGood(i)
	return ok, nil
}

func init() {
	voucher.RegisterCheckFactory("examplecheck", func() voucher.Check {
		return new(check)
	})
}
```

## Implement the Check interface

Below is the Check interface from the source code.

```golang
type Check interface {
	Check(context.Context, ImageData) (bool, error)
}
```

As you can see, Check has only one method, also called Check.

This is the method where the test itself should run.

```golang
// Check if the image described by ImageData is acceptable for our purposes.
func (c *check) Check(ctx context.Context, i voucher.ImageData) (bool, error) {
    ok := isImageGood(i)
	return ok, nil
}
```

In our example, the function `isImageGood()` is called, and it returns a boolean.

If you have a Check which might fail with an error, you should return `false, err` rather than just returning the error and not worrying about the boolean value.

### AuthorizedChecks

Some checks require authenticating against the Docker registry. These checks should implement AuthorizedCheck, which requires a new method:

```golang
type AuthCheck interface {
	Check
	SetAuth(Auth)
}
```

For example, with a check with an `auth` field, you could implement:

```golang
func (c *check) SetAuth(auth voucher.Auth) {
    c.auth = auth
}
```

Checks that implement AuthorizedCheck will automatically have their `SetAuth` method called with a usable Auth. AuthorizedChecks can also be MetadataChecks.

### MetadataChecks

Some checks require interacting with the MetadataClient. These checks should implement MetadataCheck, which requires a new method:

```golang
type MetadataCheck interface {
	Check
	SetMetadataClient(MetadataClient)
}
```

For example, if your check has a `metadataClient` field, you could write:

```golang
func (c *check) SetMetadataClient(client voucher.MetadataClient) {
    c.metadataClient = client
}
```

Checks that implement MetadataCheck will automatically have their `SetMetadataClient` method called with the configured MetadataClient. MetadataChecks can also be AuthorizedChecks.

### RepositoryChecks

Some checks require interacting with the code repository client. These checks should implement RepositoryCheck, which requires a new method:

```golang
type RepositoryCheck interface {
	MetadataCheck
	SetRepositoryClient(repository.Client)
}
```

For example, if your check has a `repositoryClient` field, you could write:

```golang
func (c *check) SetRepositoryClient(client repository.Client) {
    c.repositoryClient = client
}
```

Checks that implement RepositoryCheck will automatically have their `SetRepositoryClient` method called with the configured repository client.

## Implement the CheckFactory

CheckFactories are functions which return a new Check.

```golang
type CheckFactory func() Check
```

For simplicity's sake, in our example we are using a closure instead of a separate function.

```golang
func() voucher.Check {
    return new(check)
}
```

This method then gets passed to the `RegisterCheck` function, with the check name as the first parameter.

```golang
func init() {
	voucher.RegisterCheckFactory("examplecheck", func() voucher.Check {
		return new(check)
	})
}
```

We call `RegisterCheck` in the `init` function because it saves us from directly interacting with any Check code in the  main source tree. `RegisterCheck` adds the CheckFactory to the list of available checks. 

You can then instantiate a new "examplecheck" by calling `GetCheckFactories("examplecheck")`

```golang
checks := voucher.GetCheckFactories("examplecheck")
```

You can then execute the check by calling:

```golang
ok, err := checks["examplecheck"].Check(imageData)
```

Note that if the name given in `RegisterCheck` differs, you will access the check with that name, rather than the package name.

The final step is to Register the new check in Voucher, by adding an import line in `cmd/config/checks.go`.

```golang
import _ "github.com/grafeas/voucher/checks/examplecheck"
```

This will cause the CheckFactory to be registered.
