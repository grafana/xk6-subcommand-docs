
# unroute(url)

The method removes all routes for the `url`, that were previously added with `page.route`.

| Parameter | Type                                                                                         | Default | Description                                         |
| --------- | -------------------------------------------------------------------------------------------- | ------- | --------------------------------------------------- |
| url       | string or Regexp                                                                             | `''`    | The exact URL that was used in the `page.route(url, handler)`  call.                        |

## Returns

| Type            | Description                                                  |
| --------------- | ------------------------------------------------------------ |
| `Promise<void>` | A Promise that fulfills when the route(s) is(are) removed from the page. |

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

  await page.goto('https://quickpizza.grafana.com/');

  await page.getByRole('button', { name: 'Pizza, Please!' }).click();

  // Stop aborting all image requests
  await page.unroute(/(\.png$)|(\.jpg$)/)

  await page.close();
}
```

