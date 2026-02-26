
# Response.html()

Parses response as HTML and populate a Selection object.

### Returns

| Type                                                                                   | Description        |
| -------------------------------------------------------------------------------------- | ------------------ |
| Selection | A Selection object |

### Example

```javascript
import http from 'k6/http';

export default function () {
  const res = http.get('https://stackoverflow.com');

  const doc = res.html();
  doc
    .find('link')
    .toArray()
    .forEach(function (item) {
      console.log(item.attr('href'));
    });
}
```

