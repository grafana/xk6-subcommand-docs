
# request()

Returns the matching Request object.

## Returns

| Type      | Description           |
| --------- | --------------------- |
| `Request` | The `Request` object. |

## Example

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
    await page.route(/.*\/api\/pizza/, async function (route) {
      await route.continue({
        headers: {
          ...route.request().headers(),
          'Content-Type': 'application/json',
        },
        postData: JSON.stringify({
          customName: 'My Pizza',
        }),
      });
    });

    await page.goto('https://quickpizza.grafana.com/');
  } finally {
    await page.close();
  }
}
```

