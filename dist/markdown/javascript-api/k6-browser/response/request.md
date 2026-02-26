
# request()

Returns the matching Request object.

### Returns

| Type      | Description           |
| --------- | --------------------- |
| `Request` | The `Request` object. |

### Example

```javascript
import { browser } from 'k6/browser';

export const options = {
  scenarios: {
    ui: {
      executor: 'shared-iterations',
      options: {
        browser: {
          type: 'chromium',
        },
      },
    },
  },
};

export default async function () {
  const page = await browser.newPage();

  try {
    const res = await page.goto('https://test.k6.io/');

    const req = res.request();
    // ...
  } finally {
    await page.close();
  }
}
```

