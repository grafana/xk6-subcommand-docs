
# Worker

> **Caution:** This API is a work in progress. Some of the following functionalities might behave unexpectedly.

Represents a Web Worker or a Service Worker within the browser context.

## Supported APIs

| k6 Class                                                                               | Description                        |
| -------------------------------------------------------------------------------------- | ---------------------------------- |
| url() | Returns the URL of the web worker. |

```javascript
import { browser } from 'k6/browser';

export const options = {
  scenarios: {
    browser: {
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
    await page.goto('https://test.k6.io/browser.php');
    console.log(page.workers());
  } finally {
    await page.close();
  }
}
```

