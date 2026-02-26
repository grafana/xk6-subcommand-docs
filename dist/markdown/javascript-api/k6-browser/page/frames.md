
# frames()

Returns an array of Frames on the page.

### Returns

| Type                                                                                               | Description                                    |
| -------------------------------------------------------------------------------------------------- | ---------------------------------------------- |
| Frames[] | An array of `Frames` associated with the page. |

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
  console.log(page.frames());
}
```

