# rest-stub

Simple web application which automatically responds to HTTP requests.

Often when creating applications which consume REST API's, the resource endpoints don't yet exist, are behind authentication, or are somehow unavailable.  With this application, those requests can be pointed somewhere temporarily and provide real responses without the hassle of writing stubs.  Although, stubs are also supported.

## How It Works

By default, _all_ requests return a `200` response with no content.  Responses may be altered by the following custom headers:

- **`x-stub-delay`** - Delays response by given integer as milliseconds.
- **`x-stub-status`** - Response will return with givn status code (1xx - 5xx).
- **`x-stub-content`** - Response will return with content of specified content stub.
- - **`x-stub-content-type`** - Response will return with given `Content-Type`.

## Content Stubs

Content stubs are just text files containing data.  They are stored in a configurable directory and are addressable via the `x-stub-content` header by their file name (with extension).

For example, here is a JSON response stub located in the default stub content directory at `/var/tmp/rest-stub/example.json`:

```
HTTP/1.1 200 OK
Content-Type: application/octet-stream
Date: Mon, 25 May 2015 15:23:54 GMT
Transfer-Encoding: chunked

{
  "animal":  "duck",
  "mineral": "quartz"
}
```

Note that the `Content-Type` is `text/plain`.  You must manually define the responses `Content-Type` via the `x-stub-content-type` header or it will default to `application/octet-stream`.

## Usage

```
go install github.com/mikattack
rest-stub
```

The application supports the following options and defaults:

- **`--content`** - Stub content directory (`/var/tmp/rest-stub`)
- **`--log-level`** - Set logging level (`INFO`)
- **`--port`** - Port (`48200`)
