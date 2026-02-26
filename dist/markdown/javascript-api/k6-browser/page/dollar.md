
# page.$(selector)

> **Warning:** Use locator-based `page.locator(selector)` instead.

The method finds an element matching the specified selector within the page. If no elements match the selector, the return value resolves to `null`. To wait for an element on the page, use locator.waitFor([options]).

### Returns

| Type                             | Description                                                                                                                                                                                     |
| -------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `Promise` | A Promise that fulfills with the ElementHandle of the selector when a matching element is found or `null`. |

### Example

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

  await page.goto('https://test.k6.io/browser.php');
  const text = await page.$('#text1').then((text) => text.type('hello world'));
}
```

