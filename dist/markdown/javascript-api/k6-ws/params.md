
# Params

> **Note:** A module with a better and standard API exists.
> 
> 
> The new [k6/experimental/websockets API](/docs/k6/v1.5.0/javascript-api/k6-experimental/websockets/) partially implements the [WebSockets API living standard](https://websockets.spec.whatwg.org/).
> 
> 
> When possible, we recommend using the new API. It uses a global event loop for consistency with other k6 APIs and better performance.

_Params_ is an object used by the WebSocket methods that generate WebSocket requests. _Params_ contains request-specific options like headers that should be inserted into the request.

| Name                 | Type                                                                                        | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |
| -------------------- | ------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `Params.compression` | string                                                                                      | Compression algorithm to be used by the WebSocket connection. The only supported algorithm currently is `deflate`. If the option is left unset or empty, it defaults to no compression.                                                                                                                                                                                                                                                                                                                                                         |
| `Params.jar`         | http.CookieJar | The cookie jar that will be used when making the initial HTTP request to establish the WebSocket connection. If empty, the default VU cookie jar will be used.                                                                                                                                                                                                                                                                                                     |
| `Params.headers`     | object                                                                                      | Custom HTTP headers in key-value pairs that will be added to the initial HTTP request to establish the WebSocket connection. Keys are header names and values are header values.                                                                                                                                                                                                                                                                                                                                                                |
| `Params.tags`        | object                                                                                      | Custom metric tags in key-value pairs where the keys are names of tags and the values are tag values. The WebSocket connection will generate metrics samples with these tags attached, allowing users to filter the results data or set thresholds on sub-metrics. |

### Example of custom metadata headers and tags

_A k6 script that will make a WebSocket request with a custom header and tag results data with a specific tag_

```javascript
import ws from 'k6/ws';

export default function () {
  const url = 'wss://echo.websocket.org';
  const params = {
    headers: { 'X-MyHeader': 'k6test' },
    tags: { k6test: 'yes' },
  };
  const res = ws.connect(url, params, function (socket) {
    socket.on('open', function () {
      console.log('WebSocket connection established!');
      socket.close();
    });
  });
}
```

