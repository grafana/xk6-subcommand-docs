
# request(method, url, [body], [params])

Generic method for making arbitrary HTTP requests.

Consider using specific methods for making common requests get, post, put, patch.

| Parameter         | Type                                                                                                                              | Description                                                                                                                |
| ----------------- | --------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------- |
| method            | string                                                                                                                            | HTTP method. Must be uppercase (GET, POST, PUT, PATCH, OPTIONS, HEAD, etc)                                                 |
| url               | string                                                                                                                            | HTTP URL. If baseURL is set, provide only path.                                                                            |
| body (optional)   | null / string / object / ArrayBuffer / SharedArray | Request body; objects are `x-www-form-urlencoded`. To omit body, set to `null`.                                            |
| params (optional) | null or object {}                                                                                                                 | Additional parameters for this specific request. |

### Returns

| Type                                                                                 | Description                                                                                       |
| ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------- |
| Response | HTTP Response object. |

### Example

```javascript
import { Httpx } from 'https://jslib.k6.io/httpx/0.1.0/index.js';

const session = new Httpx({
  baseURL: 'https://quickpizza.grafana.com/api',
  timeout: 20000, // 20s timeout.
});

export default function testSuite() {
  const resp_get = session.request('GET', `/status/200`);
  const resp_post = session.request('POST', `/status/200`, { key: 'value' });
  const resp_put = session.request('PUT', `/status/200`, { key: 'value' });
  const resp_patch = session.request('PATCH', `/status/200`, { key: 'value' });
  const resp_delete = session.request('DELETE', `/status/200`);
}
```

