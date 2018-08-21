# Voucher Server API Reference

This is a quick guide to the Voucher API, for developers who are looking at building in support for Voucher.

## API Calls

### POST /all

Run all of the enabled tests on the image referred to by the passed input.

Depending on the server's configuration, your client may need to use Basic Authentication to access this call.

This call accepts a JSON encoded object with the following fields:

| Field       | Comment                                                                                     |
| :---------- | :------------------------------------------------------------------------------------------ |
| `image_url` | The URL of the image to test against. This should include the digest at the end of the URL. |

For example:

```json
{
   "image_url": "gcr.io/path/to/image@sha256:hashvalue",
}
```

The response will be a JSON encoded object containing the same fields, as well as a listing of the tests that ran, if they
were successful or not, and any errors returned during the course of the test execution.

The response will have the following fields:

| Field       | Comment                                                        |
| :---------- | :---------------------------------------------------------     |
| `image`     | The URL of the image to test against.                          |
| `success`   | A boolean, true if all tests passed, false if anyh failed.     |
| `results`   | An array of objects, with one for each test that was executed. |

The each of the objects in the `results` array are structured as follows:

| Field       | Comment                                                                            |
| :---------- | :--------------------------------------------------------------------------------- |
| `name`      | The name of the test.                                                              |
| `success`   | A boolean, true if all tests passed, false if any of the tests failed.             |
| `attested`  | A boolean, true if an attestation was created for the check.                       |
| `err`       | Any error message or structure that was thrown during the course of the execution. |

### POST /{test name here}

Run the test specified in the URL.

For example, if the API call is `POST /diy`, this will run the DIY test and nothing else.

The input and output of this API call is identical to that described in [`POST /all`](#post-all), 
and like that call, authorization may be handled by Basic Authentication.

### GET /services/ping

This call does nothing more than return a 200 Success status code. It is used to verify that the service is online.

No Authorization header is required.
