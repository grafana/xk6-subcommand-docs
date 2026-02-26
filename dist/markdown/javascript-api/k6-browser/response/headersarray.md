
# headersArray()

An array with all the response HTTP headers. Unlike response.allHeaders(), header names are not lower-cased. Headers with multiple entries, such as `Set-Cookie`, appear in the array multiple times.

### Returns

| Type                                             | Description                                                           |
| ------------------------------------------------ | --------------------------------------------------------------------- |
| `Promise` | A promise that resolves to an array of all the response HTTP headers. |

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

    const headersArray = await res.headersArray();
    console.log(`headersArray: ${JSON.stringify(headersArray)}`); // headersArray: [{"name":"Transfer-Encoding"...}]
  } finally {
    await page.close();
  }
}
```

