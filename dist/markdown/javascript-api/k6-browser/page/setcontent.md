
# setContent(html[, options])

Sets the supplied HTML string to the current page.

| Parameter         | Type   | Default | Description                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                    |
| ----------------- | ------ | ------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| html              | string | `''`    | HTML markup to assign to the page.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                             |
| options           | object | `null`  |                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                |
| options.timeout   | number | `30000` | Maximum operation time in milliseconds. Pass `0` to disable the timeout. The default value can be changed via the browserContext.setDefaultNavigationTimeout(timeout), browserContext.setDefaultTimeout(timeout), page.setDefaultNavigationTimeout(timeout) or page.setDefaultTimeout(timeout) methods. Setting the value to `0` will disable the timeout. |
| options.waitUntil | string | `load`  | When to consider operation to have succeeded. See Events for more details.                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                     |

### Returns

| Type            | Description                                                                |
| --------------- | -------------------------------------------------------------------------- |
| `Promise<void>` | A Promise that fulfills when the page has been set with the supplied HTML. |

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

  const htmlContent = `
    <!doctype html>
    <html>
      <head><meta charset='UTF-8'><title>Test</title></head>
      <body>Test</body>
    </html>
  `;

  await page.setContent(htmlContent);
}
```

