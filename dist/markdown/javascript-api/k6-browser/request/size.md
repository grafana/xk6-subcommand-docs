
# size()

This method returns the size (in bytes) of body and header sections of the Request.

### Returns

| Type          | Description            |
| ------------- | ---------------------- |
| Promise | Returns [Size](#size). |

### Size

| Property | Type   | Description                        |
| -------- | ------ | ---------------------------------- |
| body     | number | Size in bytes of the request body. |
| headers  | number | Size in bytes of the headers body. |

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

    const size = req.size();
    console.log(`size: ${JSON.stringify(size)}`); // size: {"headers":344,"body":0}
  } finally {
    await page.close();
  }
}
```

