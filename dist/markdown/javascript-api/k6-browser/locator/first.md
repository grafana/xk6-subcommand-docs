
# first()

Returns a Locator to the first matching element.

### Returns

| Type                                                                                   | Description                                               |
| -------------------------------------------------------------------------------------- | --------------------------------------------------------- |
| Locator | The first element `Locator` associated with the selector. |

### Example

```javascript
import { expect } from 'https://jslib.k6.io/k6-testing//index.js';
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
  await page.goto('https://quickpizza.grafana.com');

  await expect(await page.locator('p').first()).toContainText('QuickPizza');

  await page.close();
}
```

