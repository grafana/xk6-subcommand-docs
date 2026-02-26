
# JSHandle

Represents a reference to a JavaScript object within the context of a webpage. This allows you to interact with JavaScript objects directly from your script.

> **Caution:** This API is a work in progress. Some of the following functionalities might behave unexpectedly.

## Supported APIs

| Method                                                                                                                            | Description                                                                                                                                                                                   |
| --------------------------------------------------------------------------------------------------------------------------------- | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| asElement()                              | Returns either `null` or the object handle itself, if the object handle is an instance of ElementHandle. |
| dispose()                                  | Stops referencing the element handle.                                                                                                                                                         |
| evaluate(pageFunction[, arg])             | Evaluates the `pageFunction` and returns its return value.                                                                                                                                    |
| evaluateHandle(pageFunction[, arg]) | Evaluates the `pageFunction` and returns a `JSHandle`.                                                                                                                                        |
| getProperties()                      | Fetches a map with own property names of of the `JSHandle` with their values as `JSHandle` instances.                                                                                         |
| jsonValue()                              | Fetches a JSON representation of the object.                                                                                                                                                  |

### Example

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
    await page.goto('https://test.k6.io/');
    const jsHandle = await page.evaluateHandle(() => document.head);
    // ...
  } finally {
    await page.close();
  }
}
```

