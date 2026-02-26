
# contentFrame()

Returns the Frame that this element is contained in.

### Returns

| Type                     | Description                                                                                                                                                      |
| ------------------------ | ---------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `Promise` | A Promise that resolves to the Frame that this element is contained in |

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

  const element = await page.$('#text1');
  const frame = await element.contentFrame();
  console.log(frame.url());

  await page.close();
}
```

