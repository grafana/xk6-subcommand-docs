
# textContent()

> **Warning:** Use `locator.textContent([options])` instead.

Returns the element's text content.

### Returns

| Type                      | Description                                                              |
| ------------------------- | ------------------------------------------------------------------------ |
| `Promise<string \| null>` | A Promise that fulfills with the text content of the selector or `null`. |

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

  const options = await page.$('#checkbox1');
  console.log(await options.textContent());

  await page.close();
}
```

