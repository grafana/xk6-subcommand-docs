
# unrouteAll()

The method removes all routes created with `page.route`.

## Returns

| Type            | Description                                                  |
| --------------- | ------------------------------------------------------------ |
| `Promise<void>` | A Promise that fulfills when all routes are removed from the page. |

## Example

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

  // Abort all image requests
  await page.route(/(\.png$)|(\.jpg$)/, async function (route) {
    await route.abort();
  });

  // Fulfill request with the following response which
  // changes the quotes displayed on the page
  await page.route(/.*\/quotes$/, async function (route) {
    await route.fulfill({
      body: JSON.stringify({
        quotes: ['"We ❤️ pizza" - k6 team'],
      }),
    });
  });

  await page.goto('https://quickpizza.grafana.com/');

  await page.getByRole('button', { name: 'Pizza, Please!' }).click();

  // Stop aborting all image requests and fulfilling quote requests
  await page.unrouteAll()

  await page.close();
}
```

