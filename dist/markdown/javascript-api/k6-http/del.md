
# del( url, [body], [params] )

Make a DELETE request.

| Parameter                    | Type                                                                                            | Description                                                                                                                                                                                                                                     |
| ---------------------------- | ----------------------------------------------------------------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| url                          | string /HTTP URL | Request URL (e.g. `http://example.com`).                                                                                                                                                                                                        |
| body (optional, discouraged) | string / object / ArrayBuffer                                                                   | Request body; objects will be `x-www-form-urlencoded`. This is discouraged, because sending a DELETE request with a body has [no defined semantics](https://tools.ietf.org/html/rfc7231#section-4.3.5) and may cause some servers to reject it. |
| params (optional)            | object                                                                                          | Params object containing additional request parameters.                                                                                                               |

### Returns

| Type     | Description                                                                                       |
| -------- | ------------------------------------------------------------------------------------------------- |
| Response | HTTP Response object. |

### Example

```javascript
import http from 'k6/http';

const url = 'https://quickpizza.grafana.com/api/delete';

export default function () {
  const params = { headers: { 'X-MyHeader': 'k6test' } };
  http.del(url, null, params);
}
```

