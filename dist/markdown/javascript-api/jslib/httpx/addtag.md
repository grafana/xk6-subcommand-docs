
# addTag( key, value )

| Parameter | Type   | Description |
| --------- | ------ | ----------- |
| name      | string | Tag name    |
| value     | string | Tag value   |

### Example

```javascript
import { Httpx } from 'https://jslib.k6.io/httpx/0.1.0/index.js';

const session = new Httpx({ baseURL: 'https://quickpizza.grafana.com' });

session.addTag('tagName', 'tagValue');
session.addTag('AnotherTagName', 'tagValue2');

export default function () {
  session.get('/');
}
```

