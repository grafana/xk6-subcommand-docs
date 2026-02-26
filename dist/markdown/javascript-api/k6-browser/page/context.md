
# context()

Gets the BrowserContext that the page belongs to.

### Returns

| Type                                                                                                              | Description                                    |
| ----------------------------------------------------------------------------------------------------------------- | ---------------------------------------------- |
| BrowserContext | The `BrowserContext` that the page belongs to. |

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
  console.log(page.context()); // prints {"base_event_emitter":{}}
}
```

