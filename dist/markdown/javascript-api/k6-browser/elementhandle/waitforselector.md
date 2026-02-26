
# waitForSelector(selector[, options])

> **Warning:** Use `locator.waitFor([options])` instead.

Waits for the element to be present in the DOM and to be visible.

| Parameter       | Type    | Default   | Description                                                                                                                                                                                                                                                                                                                                   |
| --------------- | ------- | --------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| selector        | string  |           | A selector to query the element.                                                                                                                                                                                                                                                                                                              |
| options         | object  | `null`    | Optional settings.                                                                                                                                                                                                                                                                                                                            |
| options.state   | string  | `visible` | The state to wait for. This can be one of `visible`, `hidden`, `stable`, `enabled`, `disabled`, or `editable`.                                                                                                                                                                                                                                |
| options.strict  | boolean | `false`   | If set to `true`, the method will throw an error if the element is not found.                                                                                                                                                                                                                                                                 |
| options.timeout | number  | `30000`   | Maximum time in milliseconds. Pass `0` to disable the timeout. Default is overridden by the `setDefaultTimeout` option on BrowserContext or Page. |

### Returns

| Type                             | Description                                                                                                                                                                 |
| -------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `Promise` | A Promise that fulfills with the ElementHandle when the element is found. |

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
  await page.goto("https://test.k6.io");

  const element = await page.$(".header");
  const el = await element.waitForSelector(".title");
  // ... do something with the element

  await page.close();
}
```

