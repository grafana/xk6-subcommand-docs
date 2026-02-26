
# textContent(selector[, options])

> **Warning:** Use locator-based `locator.textContent([options])` instead.

Returns the `element.textContent`.

| Parameter       | Type    | Default | Description                                                                                                                                                                                                                                                                                                                                   |
| --------------- | ------- | ------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| selector        | string  | `''`    | A selector to search for an element. If there are multiple elements satisfying the selector, the first will be used.                                                                                                                                                                                                                          |
| options         | object  | `null`  |                                                                                                                                                                                                                                                                                                                                               |
| options.strict  | boolean | `false` | When `true`, the call requires selector to resolve to a single element. If given selector resolves to more than one element, the call throws an exception.                                                                                                                                                                                    |
| options.timeout | number  | `30000` | Maximum time in milliseconds. Pass `0` to disable the timeout. Default is overridden by the `setDefaultTimeout` option on BrowserContext or Page. |

### Returns

| Type                      | Description                                                              |
| ------------------------- | ------------------------------------------------------------------------ |
| `Promise<string \| null>` | A Promise that fulfills with the text content of the selector or `null`. |

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
  console.log(await page.textContent('#checkbox1'));
}
```

