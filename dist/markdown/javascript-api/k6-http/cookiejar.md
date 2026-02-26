
# cookieJar()

Get the active cookie jar.

| Type                                                                                   | Description         |
| -------------------------------------------------------------------------------------- | ------------------- |
| CookieJar | A CookieJar object. |

### Example

```javascript
import http from 'k6/http';

export default function () {
  const jar = http.cookieJar();
}
```

