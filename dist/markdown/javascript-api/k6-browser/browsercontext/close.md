
# close()

Close the browser context and all its pages. The browser context is unusable after this call and a new one must be created. This is typically called to cleanup before ending the test.

### Returns

| Type            | Description                                                                                                                                                                                                                                           |
| --------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `Promise<void>` | A Promise that fulfills when the browser context and all its pages have been closed. |

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
  const context = await browser.newContext();
  await context.newPage();

  await context.close();
}
```

