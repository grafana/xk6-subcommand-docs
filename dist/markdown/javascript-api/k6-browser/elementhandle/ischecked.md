
# isChecked()

> **Warning:** Use `locator.isChecked([options])` instead.

Checks to see if the `checkbox` `input` type is selected or not.

### Returns

| Type            | Description                                                                  |
| --------------- | ---------------------------------------------------------------------------- |
| `Promise<bool>` | A Promise that fulfills with `true` if the element is checked, else `false`. |

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

  const checkbox = await page.$('#checkbox1');
  const isChecked = await checkbox.isChecked();
  if (!isChecked) {
    await checkbox.check();
  }

  await page.close();
}
```

