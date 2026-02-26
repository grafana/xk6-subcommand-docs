
# setDefaultNavigationTimeout(timeout)

This setting will change the navigation timeout for the following methods:

- page.goto(url, [options])
- page.reload([options])
- page.setContent(html, [options])
- page.waitForNavigation([options])

| Parameter | Type   | Default | Description              |
| --------- | ------ | ------- | ------------------------ |
| timeout   | number |         | Timeout in milliseconds. |

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

  page.setDefaultNavigationTimeout(60000);
  await page.goto('https://test.k6.io/browser.php');
}
```

