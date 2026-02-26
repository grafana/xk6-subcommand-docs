
# isClosed()

> **Caution:** This method has **known issues**. For details, refer to [#878](https://github.com/grafana/xk6-browser/issues/878).

Checks if the page has been closed.

### Returns

| Type            | Description                                                                    |
| --------------- | ------------------------------------------------------------------------------ |
| `Promise<bool>` | A Promise that fullfils with `true` if the page has been closed, else `false`. |

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
  await page.close();

  console.log(await page.isClosed());
}
```

