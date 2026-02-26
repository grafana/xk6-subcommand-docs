
# response()

Returns the matching Response object, or `null` if the response was not received due to error.

### Returns

| Type                        | Description                                                                     |
| --------------------------- | ------------------------------------------------------------------------------- |
| `Promise` | The `Response` object, or `null` if the response was not received due to error. |

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

    const response = await req.response();
    // ...
  } finally {
    await page.close();
  }
}
```

