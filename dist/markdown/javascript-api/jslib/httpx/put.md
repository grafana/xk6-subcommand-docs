
# put(url, [body], [params])

`session.put(url, body, params)` makes a PUT request. Only the first parameter is required

| Parameter         | Type                                                                                                                              | Description                                                                                                                |
| ----------------- | --------------------------------------------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------- |
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
  const resp = session.put(`/put`, {
    first_name: 'Mr',
    last_name: 'Croco',
    username: 'my user',
    password: 'my password',
  });
  console.log(resp.status);
}
```

