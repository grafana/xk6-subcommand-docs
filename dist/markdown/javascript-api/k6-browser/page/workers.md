
# workers()

This method returns an array of the dedicated WebWorkers associated with the page.

### Returns

| Type                                                                                                    | Description                                     |
| ------------------------------------------------------------------------------------------------------- | ----------------------------------------------- |
| WebWorkers[] | Array of `WebWorkers` associated with the page. |

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
  console.log(page.workers());
}
```

