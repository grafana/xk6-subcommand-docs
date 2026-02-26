
# ownerFrame

Returns the Frame of the element.

### Returns

| Type                     | Description                                                                                                                                                                                           |
| ------------------------ | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `Promise` | A Promise that fulfills with the Frame of the element, or `null` if the element is not attached to a frame. |

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
  const frame = await element.ownerFrame();
  console.log(frame.url());

  await page.close();
}
```

