
# parseHTML( src )

Parse an HTML string and populate a Selection object.

| Parameter | Type   | Description  |
| --------- | ------ | ------------ |
| src       | string | HTML source. |

### Returns

| Type                                                                                   | Description         |
| -------------------------------------------------------------------------------------- | ------------------- |
| Selection | A Selection object. |

### Example

```javascript
import { parseHTML } from 'k6/html';
import http from 'k6/http';

export default function () {
  const res = http.get('https://k6.io');
  const doc = parseHTML(res.body); // equivalent to res.html()
  const pageTitle = doc.find('head title').text();
  const langAttr = doc.find('html').attr('lang');
}
```

