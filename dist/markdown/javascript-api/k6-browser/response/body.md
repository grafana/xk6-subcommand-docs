
# allHeaders()

Returns the response body.

### Returns

| Type                   | Description                                                   |
| ---------------------- | ------------------------------------------------------------- |
| `Promise` | A promise that resolves to the buffer with the response body. |

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

    const body = await res.body();
    // Refer to https://github.com/grafana/k6/issues/2440 on decoding the body to a string.
  } finally {
    await page.close();
  }
}
```

