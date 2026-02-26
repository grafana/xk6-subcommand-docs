
# evaluateHandle(pageFunction[, arg])

Executes JavaScript code in the page and returns the value of the `pageFunction` invocation as a JSHandle.

The only difference between `page.evaluate()` and `page.evaluateHandle()` is that `page.evaluateHandle()` returns JSHandle.

| Parameter    | Type               | Defaults | Description                                                              |
| ------------ | ------------------ | -------- | ------------------------------------------------------------------------ |
| pageFunction | function or string |          | Function to be evaluated in the page context. This can also be a string. |
| arg          | string             | `''`     | Optional argument to pass to `pageFunction`                              |

### Returns

| Type                | Description                                                                                                                                    |
| ------------------- | ---------------------------------------------------------------------------------------------------------------------------------------------- |
| `Promise<JSHandle>` | The JSHandle ) instance associated with the page. |

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
  const resultHandle = await page.evaluateHandle(() => document.body);
  console.log(resultHandle.jsonValue());
}
```

