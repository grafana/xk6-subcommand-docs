
# getProperties()

This allows you to inspect and interact with the properties of the object within the page context.

### Returns

| Type                           | Description                                                                         |
| ------------------------------ | ----------------------------------------------------------------------------------- |
| Promise> | A map with property names as keys and `JSHandle` instances for the property values. |

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
    await page.goto('https://test.k6.io/browser.php');
    const jsHandle = await page.evaluateHandle(() => {
      return { window, document };
    });

    const properties = await jsHandle.getProperties();
    console.log(properties); // {"window":{},"document":{}}
  } finally {
    await page.close();
  }
}
```

