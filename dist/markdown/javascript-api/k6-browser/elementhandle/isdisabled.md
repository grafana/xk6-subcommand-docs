
# isDisabled()

> **Warning:** Use `locator.isDisabled([options])` instead.

Checks if the element is disabled.

### Returns

| Type            | Description                                                                     |
| --------------- | ------------------------------------------------------------------------------- |
| `Promise<bool>` | A Promise that fulfills with `true` if the element is disabled, else `false`. |

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

  const text = await page.$('#input-text-disabled');
  const isDisabled = await text.isDisabled();
  if (isDisabled) {
    console.log('element is disabled');
  }

  await page.close();
}
```

