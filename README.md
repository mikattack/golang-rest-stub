# rest-stub

Simple web application which automatically responds to HTTP requests.

Often when creating applications which consume REST API's, the resource endpoints don't yet exist, are behind authentication, or are somehow unavailable.  With this application, those requests can be pointed somewhere temporarily and provide real responses without the hassle of writing stubs.  Although, stubs are also supported.

## How It Works

By default, _all_ requests return a `200` response with no content.  Responses may be altered by the following custom headers:

- **`x-stub-delay`** - Delays response by given integer as milliseconds.
- **`x-stub-status`** - Response will return with givn status code (1xx - 5xx).
- **`x-stub-content`** - Response will return with content of specified content stub.

## Content Stubs

Content stubs are just text files containing data.  They are stored in a configurable directory and are addressable via the `x-stub-content` header by their file name (sans-extension).

The MIME type of the content can be set by adding a single string to the top of the file, prefixed by a `#`.  Otherwise all stub content is assumed to be `application/octet-stream`.

For example, here is a JSON response stub located in the default stub content directory at `/var/tmp/rest-stub/example.json`:

```
# application/json
{
  "animal":  "duck",
  "mineral": "quartz"
}
```

## Usage

```
go install github.com/mikattack
rest-stub
```

The application supports the following options and defaults:

- **`-d`** - Stub content directory (`/var/tmp/rest-stub`)
- **`-p`** - Port (`48200`)
