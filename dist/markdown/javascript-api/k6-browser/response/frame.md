
# frame()

The Frame that initiated the request which this response is associated to.

### Returns

| Type                                                                                            | Description                                                                  |
| ----------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------- |
| Frame | The `Frame` that initiated the request which this response is associated to. |

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

    const frame = res.frame();
    // ...
  } finally {
    await page.close();
  }
}
```

