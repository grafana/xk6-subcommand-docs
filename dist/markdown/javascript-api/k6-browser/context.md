
# context()

Returns the current BrowserContext.

> **Note:** A 1-to-1 mapping between Browser and `BrowserContext` means you cannot run `BrowserContexts` concurrently. If you wish to create a new `BrowserContext` while one already exists, you will need to close the current one, and create a new one with either newContext or newPage. All resources associated to the closed `BrowserContext` will also be closed and cleaned up (such as Pages).

### Returns

| Type           | Description                                                                                                                                                 |
| -------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------- |
| object \| null | The current BrowserContext if one has been created, otherwise `null`. |

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
  console.log(browser.context()); // null

  const page1 = await browser.newPage(); // implicitly creates a new browserContext
  const context = browser.context(); // underlying live browserContext associated with browser
  const page2 = await context.newPage(); // shares the browserContext with page1

  page1.close();
  page2.close();
  await context.close();
}
```
