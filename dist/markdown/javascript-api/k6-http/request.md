
# request( method, url, [body], [params] )

| Parameter         | Type                                                                                            | Description                                                                                                                       |
| ----------------- | ----------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------- |
| method            | string                                                                                          | Request method (e.g. `'POST'`). Must be uppercase.                                                                                |
| url               | string /HTTP URL | Request URL (e.g. `'http://example.com'`).                                                                                        |
| body (optional)   | string / object / ArrayBuffer                                                                   | Request body; Objects will be `x-www-form-urlencoded` encoded.                                                                    |
| params (optional) | object                                                                                          | Params object containing additional request parameters. |

### Returns

| Type     | Description                                                                                       |
| -------- | ------------------------------------------------------------------------------------------------- |
| Response | HTTP Response object. |

### Example

Using http.request() to issue a POST request:

```javascript
import http from 'k6/http';

const url = 'https://quickpizza.grafana.com/api/post';

export default function () {
  const data = { name: 'Bert' };

  // Using a JSON string as body
  let res = http.request('POST', url, JSON.stringify(data), {
    headers: { 'Content-Type': 'application/json' },
  });
  console.log(res.json().name); // Bert

  // Using an object as body, the headers will automatically include
  // 'Content-Type: application/x-www-form-urlencoded'.
  res = http.request('POST', url, data);
  console.log(res.body); // name=Bert
}
```

