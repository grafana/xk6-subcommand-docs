
# get( url, [params] )

Make a GET request.

| Parameter         | Type                                                                                            | Description                                                                                                                       |
| ----------------- | ----------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------- |
| url               | string /HTTP URL | Request URL (e.g. `http://example.com`).                                                                                          |
| params (optional) | object                                                                                          | Params object containing additional request parameters. |

### Returns

| Type                                                                                 | Description           |
| ------------------------------------------------------------------------------------ | --------------------- |
| Response | HTTP Response object. |

### Example fetching a URL

```javascript
import http from 'k6/http';

export default function () {
  const res = http.get('https://test.k6.io');
  console.log(JSON.stringify(res.headers));
}
```

