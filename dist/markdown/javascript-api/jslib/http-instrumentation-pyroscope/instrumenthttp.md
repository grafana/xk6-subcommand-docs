
# instrumentHTTP

The `instrumentHTTP` function instruments the k6 HTTP module with additional headers. It transparently replaces each of the k6 http module functions with versions that automatically attach a baggage header(by default) to every request.

Instrumented functions include del,get,head,options,patch,post,put, and request.

This means that, to instrument the HTTP requests, you don't need to rewrite any code.
Instead, call the module once in the init context.
For details about propagation, refer to About baggage header.

## Parameters

| Name              | Type                                        | Description                                                                                                                                                                                                                                                                                                    |
| :---------------- | :------------------------------------------ | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `generateHeaders` | `function(method, body, params) headersMap` | A function taking the request method, body and params and returning a map of headers that will be attached. By default adding baggage header as explained in the about baggage header |

## Example

This example demonstrates how to use the this library to instrument every HTTP request made in a script with baggage header.

```javascript
import { check } from 'k6';
import pyroscope from 'https://jslib.k6.io/http-instrumentation-pyroscope//index.js';
import http from 'k6/http';

// instrumentHTTP will ensure that all requests made by the `http` module
// from this point forward will have a baggage context attached.
pyroscope.instrumentHTTP();

export default () => {
  // the instrumentHTTP call in the init context replaced
  // the `http` module with a version that will automatically
  // attach a baggage header to every request.
  const res = http.get('http://httpbin.org/get', {
    headers: {
      'X-Example-Header': 'instrumented/get',
    },
  });
};
```

