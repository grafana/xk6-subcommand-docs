
# evaluateHandle(pageFunction[, arg])

Executes JavaScript code in the page, passing the matching element of the locator as the first argument to the `pageFunction` and `arg` as the following arguments. It returns the value of the `pageFunction` invocation as a JSHandle.

The only difference between `evaluate` and `evaluateHandle` is that `evaluateHandle` returns JSHandle.

| Parameter    | Type               | Defaults | Description                                  |
| ------------ | ------------------ | -------- | -------------------------------------------- |
| pageFunction | function or string |          | Function to be evaluated in the page context.                    |
| arg          | string             | `''`     | Optional argument to pass to `pageFunction`. |

### Returns

| Type              | Description                                         |
| ----------------- | --------------------------------------------------- |
| Promise<JSHandle> | A [JSHandle]((https://grafana.com/docs/k6/v1.5.0/javascript-api/k6-browser/jshandle/)) of the return value of `pageFunction`. |

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

  try {
    await page.goto("https://quickpizza.grafana.com", { waitUntil: "load" });

    await page.getByText('Pizza, Please!').click();

    const jsHandle = await page.locator('#pizza-name').evaluateHandle((pizzaName) => pizzaName);

    const obj = await jsHandle.evaluateHandle((handle) => {
      return { innerText: handle.innerText };
    });
    console.log(await obj.jsonValue()); // {"innerText":"Our recommendation:"}
  } finally {
    await page.close();
  }
}
```

