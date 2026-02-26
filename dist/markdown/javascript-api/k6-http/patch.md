
# patch( url, [body], [params] )

| Parameter         | Type                                                                                            | Description                                                                                                                      |
| ----------------- | ----------------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------------------------------------------------------- |
| url               | string /HTTP URL | Request URL (e.g. `http://example.com`).                                                                                         |
| body (optional)   | string / object / ArrayBuffer                                                                   | Request body; objects will be `x-www-form-urlencoded`.                                                                           |
| params (optional) | object                                                                                          | Params object containing additional request parameters |

### Returns

| Type                                                                                 | Description           |
| ------------------------------------------------------------------------------------ | --------------------- |
| Response | HTTP Response object. |

### Example

```javascript
import http from 'k6/http';

const url = 'https://quickpizza.grafana.com/api/patch';

export default function () {
  const headers = { 'Content-Type': 'application/json' };
  const data = { name: 'Bert' };

  const res = http.patch(url, JSON.stringify(data), { headers: headers });

  console.log(JSON.parse(res.body).name);
}
```

