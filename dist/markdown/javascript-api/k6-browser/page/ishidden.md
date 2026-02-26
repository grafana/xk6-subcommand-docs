
# isHidden(selector[, options])

> **Warning:** Use locator-based [`locator.isHidden([options])` instead.

Checks if the element is `hidden`.

| Parameter      | Type    | Default | Description                                                                                                                                                        |
| -------------- | ------- | ------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| selector       | string  | `''`    | A selector to search for an element. If there are multiple elements satisfying the selector, the first one will be used.                                           |
| options        | object  | `null`  |                                                                                                                                                                    |
| options.strict | boolean | `false` | When `true`, the call requires the selector to resolve to a single element. If the given selector resolves to more than one element, the call throws an exception. |

### Returns

| Type            | Description                                                                   |
| --------------- | ----------------------------------------------------------------------------- |
| `Promise<bool>` | A Promise that fullfils with `true` if the element is `hidden`, else `false`. |

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
  const isHidden = await page.isHidden('#input-text-hidden');
  if (isHidden) {
    console.log('element is hidden');
  }
}
```

