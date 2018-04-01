# rest-stub

Simple web application which automatically responds to HTTP requests.

Often when creating applications which consume REST API's, the resource endpoints don't yet exist, are behind authentication, or are somehow unavailable.  With this application, those requests can be pointed somewhere temporarily and get real responses.

Think of it as a useful tool until you have a test environment set up.

## How It Works

By default, _all_ requests return a `200` response with no content.  Responses may be altered by including any of the following custom headers in a request:

- **`x-stub-content-type`** - Change the value of response's `Content-Type` character set.
- **`x-stub-content`** - Response body will be the contents of a file within the configured content directory (file name is the header value). These are referred to as **content stubs**.
- **`x-stub-content-type`** - Sets the `Content-Type` of the response.
- **`x-stub-delay`** - Delays the response by a given number of milliseconds.
- **`x-stub-status`** - Sets the response's status code (1xx - 5xx).

## Content Stubs

Content stubs are just text files containing data.  They are stored in a configurable directory and are addressable via the `x-stub-content` header by their file name (with extension).

For example, here is a JSON response stub located in the default stub content directory at `/var/tmp/rest-stub/example.json`:

```
{
  "animal":     "duck",
  "mineral":    "quartz",
  "vegetable":  "kale"
}
```

...will output:

```
HTTP/1.1 200 OK
Content-Length: 77
Content-Type: text/plain; charset=utf-8
Date: Sun, 01 Apr 2018 18:35:36 GMT

{
    "animal": "duck",
    "mineral": "quartz",
    "vegetable": "kale"
}
```

Note that the `Content-Type` is `text/plain`.  You must manually define the responses `Content-Type` via the `x-stub-content-type` header to make it something more JSON-like, as content stub processing is not MIME aware.

## Usage

```
go install https://github.com/mikattack/golang-rest-stub
make server
rest-stub
```

The application supports the following options and defaults:

- **`--content`** - Stub content directory (`/var/tmp/rest-stub`)
- **`--port`** - Port (`48200`)
