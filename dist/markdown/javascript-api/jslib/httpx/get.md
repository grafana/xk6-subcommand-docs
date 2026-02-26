
# get(url, [body], [params])

`session.get(url, body, params)` makes a GET request. Only the URL parameter is required

| Parameter         | Type                                                                                                                              | Description                                                                                                                |
| ----------------- | --------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------- |
| url               | string                                                                                                                            | HTTP URL. If baseURL is set, provide only path.                                                                            |
| body (optional)   | null / string / object / ArrayBuffer / SharedArray | Set to `null` to omit the body.                                                                                            |
| params (optional) | null or object {}                                                                                                                 | Additional parameters for this specific request. |

### Returns

| Type                                                                                 | Description                                                                                       |
| ------------------------------------------------------------------------------------ | ------------------------------------------------------------------------------------------------- |
| Response | HTTP Response object. |

### Example

```javascript
import { Httpx } from 'https://jslib.k6.io/httpx/0.1.0/index.js';

const session = new Httpx({
  baseURL: 'https://quickpizza.grafana.com',
  timeout: 20000, // 20s timeout.
});

export default function testSuite() {
  const resp = session.get(`/api/names`);
}
```

