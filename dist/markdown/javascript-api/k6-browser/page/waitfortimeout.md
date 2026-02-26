
# waitForTimeout(timeout)

> **Note:** Never wait for timeout in production, use this only for debugging. Tests that wait for time are inherently flaky. Use `Locator` actions and web assertions that wait automatically.

Waits for the given `timeout` in milliseconds.

### Returns

| Type            | Description                                          |
| --------------- | ---------------------------------------------------- |
| `Promise<void>` | A Promise that fulfills when the timeout is reached. |

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
  await page.waitForTimeout(5000);
}
```

