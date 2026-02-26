
# options( url, [body], [params] )

| Parameter         | Type                                                                                            | Description                                                                                                                       |
| ----------------- | ----------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------- |
| url               | string /HTTP URL | Request URL (e.g. `http://example.com`).                                                                                          |
| body (optional)   | string / object / ArrayBuffer                                                                   | Request body; objects will be `x-www-form-urlencoded`.                                                                            |
| params (optional) | object                                                                                          | Params object containing additional request parameters. |

### Returns

| Type     | Description                                                                                       |
| -------- | ------------------------------------------------------------------------------------------------- |
| Response | HTTP Response object. |

### Example

```javascript
import http from 'k6/http';

const url = 'https://quickpizza.grafana.com/';

export default function () {
  const params = { headers: { 'X-MyHeader': 'k6test' } };
  const res = http.options(url, null, params);
  console.log(res.headers['Allow']);
}
```

